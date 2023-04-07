package gcs

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
	"github.com/padok-team/yatas-gcp/internal"
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
			fmt.Printf("Failed to list buckets: %v", err)
		}
		buckets = append(buckets, *bucketAttrs)
	}

	return buckets
}
