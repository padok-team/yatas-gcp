package main

import (
	"encoding/gob"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/stangirard/yatas/plugins/commons"
)

type YatasPlugin struct {
	logger hclog.Logger
}

// Don't remove this function
func (g *YatasPlugin) Run(c *commons.Config) []commons.Tests {
	g.logger.Debug("message from Yatas Template Plugin")
	var err error
	if err != nil {
		panic(err)
	}
	var checksAll []commons.Tests

	checks, err := runPlugin(c, "template")
	if err != nil {
		g.logger.Error("Error running plugins", "error", err)
	}
	checksAll = append(checksAll, checks...)
	return checksAll
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


func authenticateImplicitWithAdc(w io.Writer, projectId string) error {
         projectId := "padok-lab"

        ctx := context.Background()

        // NOTE: Replace the client created below with the client required for your application.
        // Note that the credentials are not specified when constructing the client.
        // The client library finds your credentials using ADC.
        client, err := storage.NewClient(ctx)
        if err != nil {
                return fmt.Errorf("NewClient: %v", err)
        }
        defer client.Close()

        it := client.Buckets(ctx, projectId)
        for {
                bucketAttrs, err := it.Next()
                if err == iterator.Done {
                        break
                }
                if err != nil {
                        return err
                }
                fmt.Fprintf(w, "Bucket: %v\n", bucketAttrs.Name)
        }

        fmt.Fprintf(w, "Listed all storage buckets.\n")

        return nil
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
		"template": &commons.YatasPlugin{Impl: yatasPlugin},
	}

	logger.Debug("message from plugin", "foo", "bar")

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}

// Function that runs the checks or things to dot
func runPlugin(c *commons.Config, plugin string) ([]commons.Tests, error) {
	var checksAll []commons.Tests

	

	return checksAll, nil
}
