package gcs

import (
	"context"
	"sync"

	"cloud.google.com/go/iam"
	"cloud.google.com/go/storage"

	"github.com/padok-team/yatas-gcp/internal"
	"github.com/padok-team/yatas-gcp/logger"
	"github.com/padok-team/yatas/plugins/commons"
)

type GCSBucket struct {
	Bucket storage.BucketAttrs
	Policy iam.Policy
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
		{
			Title:          "GCP_GCS_002",
			Description:    "Check if GCS buckets are encrypted with a custom KMS key",
			Categories:     []string{"Security", "Good Practice"},
			ConditionFn:    GCSBucketEncryptionEnabled,
			SuccessMessage: "GCS bucket is encrypted with a custom KMS key",
			FailureMessage: "GCS bucket is not encrypted",
		},
		{
			Title:          "GCP_GCS_003",
			Description:    "Check if GCS buckets are not public",
			Categories:     []string{"Security", "Good Practice"},
			ConditionFn:    GCSBucketNoPublicAccess,
			SuccessMessage: "GCS bucket are not public",
			FailureMessage: "GCS bucket are public",
		},
	}

	var resources []commons.Resource
	for _, bucket := range buckets {
		policy := GetBucketPolicy(account, client, bucket.Name)
		resources = append(resources, &GCSBucket{Bucket: bucket, Policy: policy})
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
