package gcs

import (
	"context"

	"cloud.google.com/go/iam"
	"cloud.google.com/go/storage"
	"github.com/padok-team/yatas-gcp/internal"
	"github.com/padok-team/yatas-gcp/logger"
	"google.golang.org/api/iterator"
)

func GetBuckets(account internal.GCPAccount, client *storage.Client) []storage.BucketAttrs {
	var buckets []storage.BucketAttrs

	it := client.Buckets(context.TODO(), account.Project)
	for {
		bucketAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			logger.Logger.Error("Failed to list buckets", "error", err.Error())
			break
		}
		buckets = append(buckets, *bucketAttrs)
	}
	return buckets
}

func GetBucketPolicy(account internal.GCPAccount, client *storage.Client, bucket string) iam.Policy {
	c := client.Bucket(bucket)
	policy, err := c.IAM().Policy(context.TODO())
	if err != nil {
		logger.Logger.Error("Failed to get bucket policy", "error", err.Error())
		return iam.Policy{}
	}
	return *policy
}
