package main

import (
	"encoding/gob"
	"os"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/stangirard/yatas/plugins/commons"
	"context"
	"fmt"
	"sync"
    "google.golang.org/api/storage/v1"
	"github.com/padok-team/yatas-gcp/gcp/gcs"
)


type GCP_Account struct {
	Project  string `yaml:"project"`    // Name of the account in the reports
}

type YatasPlugin struct {
	logger hclog.Logger
}

// Don't remove this function
// funcrion of Yatas plugin
func (g *YatasPlugin) Run(c *commons.Config) []commons.Tests {
	g.logger.Debug("message from Yatas Template Plugin")
	var err error
	var accounts []GCP_Account
	accounts, err = UnmarshalGCP(g, c)
	g.logger.Debug("check",accounts)
	if err != nil {
		panic(err)
	}

	var checksAll []commons.Tests

	checks, err := runPlugins(c, "gcp",accounts)
	if err != nil {
		g.logger.Error("Error running plugins", "error", err)
	}
	checksAll = append(checksAll, checks...)

    // if err != nil {
    //     g.logger.Debug("Failed to connect to GCP: %v", err)
    // }
	// buckets, err := client.Buckets.List(accounts[0].Project).Do()
    // if err != nil {
    //     g.logger.Debug("Failed to list GCP buckets: %v", err)
    // }
    // for _, bucket := range buckets.Items {
    //     g.logger.Debug("Bucket: %v\n", bucket.Name)
    // }	
	return checksAll
}


// Run the plugins that are enabled in the config with a switch based on the name of the plugin
func runPlugins(c *commons.Config, plugin string,accounts []GCP_Account) ([]commons.Tests, error) {

	var checksAll []commons.Tests

	
	checksAll, err := Run(c,accounts)
	if err != nil {
		return nil, err
	}

	return checksAll, nil
}

func Run(c *commons.Config,accounts []GCP_Account ) ([]commons.Tests, error) {

	var wg sync.WaitGroup
	var queue = make(chan commons.Tests, 10)
	var checks []commons.Tests
	client,_ := connectToGCP(accounts)
	fmt.Println("value client",client)	
	wg.Add(len(accounts))
	for _, account := range accounts {
		go runTestsForAccount(client,account, c, queue)
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

func runTestsForAccount(client *storage.Service ,account GCP_Account, c *commons.Config, queue chan commons.Tests) {
	
	checks := initTest(client, c, account)
	queue <- checks
}


func initTest(s *storage.Service, c *commons.Config, a GCP_Account) commons.Tests {

	var checks commons.Tests
	checks.Account = a.Project
	var wg sync.WaitGroup
	queue := make(chan []commons.Check, 100)
	go commons.CheckMacroTest(&wg, c, gcs.RunChecks)(&wg, s, c, queue)


	go func() {
		for t := range queue {

			checks.Checks = append(checks.Checks, t...)

			wg.Done()

		}
	}()
	wg.Wait()

	return checks
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

func connectToGCP([]GCP_Account) (*storage.Service, error) {
    ctx := context.Background()
    // clientOptions := []option.ClientOption{}
    service, err := storage.NewService(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to create GCP storage client: %v", err)
    }
    return service, nil
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

func UnmarshalGCP(g *YatasPlugin, c *commons.Config) ([] GCP_Account, error) {
	var accounts [] GCP_Account

	for _, r := range c.PluginConfig {
		var tmpAccounts [] GCP_Account
		gcpFound := false
		for key, value := range r {

			switch key {
			case "pluginName":
				if value == "gcp" {
					gcpFound = true

				}
			case "accounts":

				for _, v := range value.([]interface{}) {
					var account GCP_Account
					g.logger.Debug("🔎")
					g.logger.Debug("%v", v)
					for keyaccounts, valueaccounts := range v.(map[string]interface{}) {
						switch keyaccounts {
						case "project":
							account.Project = valueaccounts.(string)
						}
					}
					tmpAccounts = append(tmpAccounts, account)

				}

			}
		}
		if gcpFound {
			g.logger.Debug("✅✅")
			accounts = tmpAccounts
		}

	}
	g.logger.Debug("✅")
	g.logger.Debug("%v", accounts)
	g.logger.Debug("Length of accounts: %d", len(accounts))
	g.logger.Debug("test",accounts[0].Project)
	return accounts, nil
}
