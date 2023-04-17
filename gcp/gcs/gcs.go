package gcs

import (
	"context"
	"sync"

	"cloud.google.com/go/storage"

	"github.com/padok-team/yatas-gcp/internal"
	"github.com/padok-team/yatas-gcp/logger"
	"github.com/padok-team/yatas/plugins/commons"
)

type GCSBucket struct {
	Bucket storage.BucketAttrs
}

func (b *GCSBucket) GetID() string {
	return b.Bucket.Name
}

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

	gcsChecks := []commons.CheckDefinition{
		{
			Title:          "GCP_GCS_001",
			Description:    "Check if GCS buckets are using object versioning",
			Categories:     []string{"Security", "Good Practice"},
			ConditionFn:    GCSBucketVersioningEnabled,
			SuccessMessage: "GCS bucket is using object versioning",
			FailureMessage: "GCS bucket is not using object versioning",
		},
	}

	var resources []commons.Resource
	for _, bucket := range buckets {
		resources = append(resources, &GCSBucket{Bucket: bucket})
	}
	commons.AddChecks(&checkConfig, gcsChecks)
	go commons.CheckResources(checkConfig, resources, gcsChecks)

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
