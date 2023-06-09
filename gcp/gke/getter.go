package gke

import (
	"context"

	container "cloud.google.com/go/container/apiv1"
	containerpb "cloud.google.com/go/container/apiv1/containerpb"
	"github.com/padok-team/yatas-gcp/internal"
	"github.com/padok-team/yatas-gcp/logger"
)

// GetClusters returns all the clusters of the project (all zones/regions)
func GetClusters(account internal.GCPAccount) []*containerpb.Cluster {
	ctx := context.Background()
	c, err := container.NewClusterManagerClient(ctx)
	if err != nil {
		logger.Logger.Error("Failed to create container client", "error", err)
		return nil
	}
	req := &containerpb.ListClustersRequest{
		Parent: "projects/" + account.Project + "/locations/-",
	}
	resp, err := c.ListClusters(ctx, req)
	if err != nil {
		logger.Logger.Error("Failed to list clusters", "error", err)
		return nil
	}
	return resp.Clusters
}
