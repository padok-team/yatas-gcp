package gcs

import (
	"context"
	"fmt"
	"sync"

	"cloud.google.com/go/storage"
	"github.com/padok-team/yatas-gcp/internal"
	"github.com/padok-team/yatas/plugins/commons"
)

func RunChecks(wa *sync.WaitGroup, account internal.GCPAccount, c *commons.Config, queue chan []commons.Check) {
	// TODO: Use commons.CheckConfig to configure the checks, once it's modified to support GCP in Yatas

	checkQueue := make(chan commons.Check, 10)
	wg := &sync.WaitGroup{}

	var checks []commons.Check

	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		fmt.Printf("Failed to create client: %v", err)
	}
	defer client.Close()

	buckets := GetBuckets(account, client)

	go commons.CheckTest(wg, c, "GCP_GCS_001", CheckIfVersioningEnabled)(checkQueue, buckets, "GCP_GCS_001")
	// Wait for all the goroutines to finish

	go func() {
		for t := range checkQueue {
			t.EndCheck()
			checks = append(checks, t)

			wg.Done()

		}
	}()

	wg.Wait()

	queue <- checks
}
