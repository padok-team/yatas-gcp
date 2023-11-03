package cloudrun

import (
	"context"

	run "cloud.google.com/go/run/apiv2"
	"cloud.google.com/go/run/apiv2/runpb"
	"github.com/padok-team/yatas-gcp/internal"
	"github.com/padok-team/yatas-gcp/logger"
	"google.golang.org/api/iterator"
)

func GetCloudRunServices(account internal.GCPAccount) []runpb.Service {
	ctx := context.Background()

	c, err := run.NewServicesClient(ctx)
	if err != nil {
		logger.Logger.Error("Failed to create CloudRun Services client", "error", err)
	}
	defer c.Close()

	var services []runpb.Service

	for _, region := range account.ComputeRegions {
		req := &runpb.ListServicesRequest{
			Parent: "projects/" + account.Project + "/locations/" + region,
		}
		it := c.ListServices(ctx, req)
		for {
			resp, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				logger.Logger.Error("Failed to list CloudRun services", "error", err.Error())
				break
			}
			services = append(services, *resp)
		}
	}

	return services

}
