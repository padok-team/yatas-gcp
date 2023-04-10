package gcs

import (
	"context"
	"sync"

	"cloud.google.com/go/storage"
	"github.com/padok-team/yatas-gcp/internal"
	"github.com/padok-team/yatas-gcp/logger"
	"github.com/padok-team/yatas/plugins/commons"
)

func RunChecks(wa *sync.WaitGroup, account internal.GCPAccount, c *commons.Config, queue chan []commons.Check) {
	var checkConfig commons.CheckConfig
	checkConfig.Init(c)

	var checks []commons.Check

	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		logger.Logger.Error("Failed to create storage client", "error", err)
	}
	defer client.Close()

	buckets := GetBuckets(account, client)

	go commons.CheckTest(checkConfig.Wg, c, "GCP_GCS_001", CheckIfVersioningEnabled)(checkConfig.Queue, buckets, "GCP_GCS_001")
	// Wait for all the goroutines to finish

	go func() {
		for t := range checkConfig.Queue {
			t.EndCheck()
			checks = append(checks, t)

			checkConfig.Wg.Done()
		}
	}()

	checkConfig.Wg.Wait()

	queue <- checks
}
