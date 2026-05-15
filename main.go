package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"os"
	"sync"

	"cloud.google.com/go/resourcemanager/apiv3"
	resourcemanagerpb "cloud.google.com/go/resourcemanager/apiv3/resourcemanagerpb"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/padok-team/yatas-gcp/gcp/cloudrun"
	"github.com/padok-team/yatas-gcp/gcp/functions"
	"github.com/padok-team/yatas-gcp/gcp/gcs"
	"github.com/padok-team/yatas-gcp/gcp/gke"
	"github.com/padok-team/yatas-gcp/gcp/iam"
	"github.com/padok-team/yatas-gcp/gcp/instance"
	"github.com/padok-team/yatas-gcp/gcp/loadbalancing"
	"github.com/padok-team/yatas-gcp/gcp/network"
	"github.com/padok-team/yatas-gcp/gcp/sql"
	"github.com/padok-team/yatas-gcp/internal"
	"github.com/padok-team/yatas-gcp/logger"
	"github.com/padok-team/yatas/plugins/commons"
	"google.golang.org/api/iterator"
)

type YatasPlugin struct {
	logger hclog.Logger
}

// Don't remove this function
// function of Yatas plugin
func (g *YatasPlugin) Run(c *commons.Config) []commons.Tests {
	logger.Logger = g.logger
	var err error
	var accounts []internal.GCPAccount
	accounts, err = UnmarshalGCP(g, c)
	if err != nil {
		logger.Logger.Error("Error unmarshaling GCP accounts", "error", err)
		return nil
	}

	var checksAll []commons.Tests

	checks, err := runPlugins(c, "gcp", accounts)
	if err != nil {
		logger.Logger.Error("Error running plugins", "error", err)
	}
	checksAll = append(checksAll, checks...)

	return checksAll
}

// Run the plugins that are enabled in the config with a switch based on the name of the plugin
func runPlugins(c *commons.Config, plugin string, accounts []internal.GCPAccount) ([]commons.Tests, error) {

	var checksAll []commons.Tests

	checksAll, err := Run(c, accounts)
	if err != nil {
		return nil, err
	}

	return checksAll, nil
}

func Run(c *commons.Config, accounts []internal.GCPAccount) ([]commons.Tests, error) {

	var wg sync.WaitGroup
	var queue = make(chan commons.Tests, 10)
	var checks []commons.Tests

	wg.Add(len(accounts))
	for _, account := range accounts {
		go runTestsForAccount(account, c, queue)
	}
	go func() {
		for t := range queue {
			checks = append(checks, t)

			wg.Done()
		}
	}()
	wg.Wait()

	return checks, nil
}

func runTestsForAccount(account internal.GCPAccount, c *commons.Config, queue chan commons.Tests) {
	checks := initTest(account, c)
	queue <- checks
}

func initTest(account internal.GCPAccount, c *commons.Config) commons.Tests {
	var checks commons.Tests
	checks.Account = account.Project
	var wg sync.WaitGroup
	queue := make(chan []commons.Check, 100)

	go commons.CheckMacroTest(&wg, c, gcs.RunChecks)(&wg, account, c, queue)
	go commons.CheckMacroTest(&wg, c, instance.RunChecks)(&wg, account, c, queue)
	go commons.CheckMacroTest(&wg, c, sql.RunChecks)(&wg, account, c, queue)
	go commons.CheckMacroTest(&wg, c, loadbalancing.RunChecks)(&wg, account, c, queue)
	go commons.CheckMacroTest(&wg, c, gke.RunChecks)(&wg, account, c, queue)
	go commons.CheckMacroTest(&wg, c, cloudrun.RunChecks)(&wg, account, c, queue)
	go commons.CheckMacroTest(&wg, c, functions.RunChecks)(&wg, account, c, queue)
	go commons.CheckMacroTest(&wg, c, iam.RunChecks)(&wg, account, c, queue)
	go commons.CheckMacroTest(&wg, c, network.RunChecks)(&wg, account, c, queue)

	go func() {
		for t := range queue {

			checks.Checks = append(checks.Checks, t...)

			wg.Done()

		}
	}()
	wg.Wait()

	return checks
}

func parseGCPAccount(values map[string]interface{}) internal.GCPAccount {
	var account internal.GCPAccount
	for keyaccounts, valueaccounts := range values {
		switch keyaccounts {
		case "project":
			account.Project = valueaccounts.(string)
		case "organization":
			account.Organization = valueaccounts.(string)
		case "folder":
			account.Folder = valueaccounts.(string)
		case "computeRegions":
			// Cannot directly unmarshal []interface{} to []string
			computeRegions := valueaccounts.([]interface{})
			for _, computeRegion := range computeRegions {
				account.ComputeRegions = append(account.ComputeRegions, computeRegion.(string))
			}
		}
	}
	return account
}

func expandGCPAccount(account internal.GCPAccount) ([]internal.GCPAccount, error) {
	if account.Project != "" {
		return []internal.GCPAccount{account}, nil
	}

	ctx := context.Background()
	var parents []string
	if account.Organization != "" {
		parents = append(parents, fmt.Sprintf("organizations/%s", account.Organization))
	}
	if account.Folder != "" {
		parents = append(parents, fmt.Sprintf("folders/%s", account.Folder))
	}

	var accounts []internal.GCPAccount
	for _, parent := range parents {
		projects, err := listProjects(ctx, parent, account.ComputeRegions)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, projects...)
		if err := listFolderProjects(ctx, parent, account.ComputeRegions, &accounts); err != nil {
			return nil, err
		}
	}
	return accounts, nil
}

func listProjects(ctx context.Context, parent string, computeRegions []string) ([]internal.GCPAccount, error) {
	client, err := resourcemanager.NewProjectsClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	it := client.ListProjects(ctx, &resourcemanagerpb.ListProjectsRequest{Parent: parent})
	var accounts []internal.GCPAccount
	for {
		project, err := it.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			return nil, err
		}
		accounts = append(accounts, internal.GCPAccount{Project: project.ProjectId, ComputeRegions: computeRegions})
	}
	return accounts, nil
}

func listFolderProjects(ctx context.Context, parent string, computeRegions []string, accounts *[]internal.GCPAccount) error {
	client, err := resourcemanager.NewFoldersClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	it := client.ListFolders(ctx, &resourcemanagerpb.ListFoldersRequest{Parent: parent})
	for {
		folder, err := it.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			return err
		}
		projects, err := listProjects(ctx, folder.Name, computeRegions)
		if err != nil {
			return err
		}
		*accounts = append(*accounts, projects...)
		if err := listFolderProjects(ctx, folder.Name, computeRegions, accounts); err != nil {
			return err
		}
	}
	return nil
}

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  2,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

func main() {
	gob.Register([]interface{}{})
	gob.Register(map[string]interface{}{})
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	yatasPlugin := &YatasPlugin{
		logger: logger,
	}
	// pluginMap is the map of plugins we can dispense.
	// Name of your plugin
	var pluginMap = map[string]plugin.Plugin{
		"gcp": &commons.YatasPlugin{Impl: yatasPlugin},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}

func UnmarshalGCP(g *YatasPlugin, c *commons.Config) ([]internal.GCPAccount, error) {
	var accounts []internal.GCPAccount

	for _, r := range c.PluginConfig {
		var tmpAccounts []internal.GCPAccount
		gcpFound := false
		for key, value := range r {

			switch key {
			case "pluginName":
				if value == "gcp" {
					gcpFound = true

				}
			case "accounts":

				for _, v := range value.([]interface{}) {
					account := parseGCPAccount(v.(map[string]interface{}))
					logger.Logger.Debug("Inspecting account", "account", v)
					resolvedAccounts, err := expandGCPAccount(account)
					if err != nil {
						return nil, err
					}
					tmpAccounts = append(tmpAccounts, resolvedAccounts...)

				}

			}
		}
		if gcpFound {
			logger.Logger.Debug("GCP config found ✅")
			accounts = tmpAccounts
		}

	}
	logger.Logger.Debug("Unmarshal Done ✅")
	logger.Logger.Debug("All accounts", "accounts", accounts)
	logger.Logger.Debug("Length of accounts", "len", len(accounts))
	return accounts, nil
}
