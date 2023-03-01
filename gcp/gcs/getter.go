package gcs

import (
	"fmt"

	"github.com/padok-team/yatas-gcp/internal"
	"google.golang.org/api/storage/v1"
)

func GetListGCS(client internal.Client_Account) *storage.Buckets {

	buckets, err := client.Client.Buckets.List(client.Gcp_account.Project).Do()
	if err != nil {
		fmt.Println("Failed to list GCP buckets: %v", err)
	}
	for _, bucket := range buckets.Items {
		fmt.Println("Bucket: %v\n", bucket.Name)
	}

	return buckets
}
