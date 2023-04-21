package instance

import (
	"context"
	"strings"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"github.com/padok-team/yatas-gcp/internal"
	"github.com/padok-team/yatas-gcp/logger"
	"google.golang.org/api/iterator"
)

// Get the compute zones based on the list of compute regions provided in the config
func GetComputeZones(account internal.GCPAccount) []string {
	ctx := context.Background()
	client, err := compute.NewZonesRESTClient(ctx)
	if err != nil {
		logger.Logger.Error("Failed to create Compute Zones client", "error", err)
	}
	defer client.Close()

	req := &computepb.ListZonesRequest{
		Project: account.Project,
	}
	var zones []string
	it := client.List(context.TODO(), req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			logger.Logger.Error("Failed to list Compute Zones", "error", err.Error())
			break
		}

		region := resp.GetRegion()
		for _, r := range account.ComputeRegions {
			if strings.HasSuffix(region, r) {
				zones = append(zones, resp.GetName())
			}
		}
	}

	logger.Logger.Debug("Compute Zones", "zones", zones)

	return zones
}

// Get all the VM instances of the account for the given compute zone
func GetInstances(account internal.GCPAccount, client *compute.InstancesClient, computeZone string) []computepb.Instance {
	req := &computepb.ListInstancesRequest{
		Project: account.Project,
		Zone:    computeZone,
	}
	var instances []computepb.Instance
	it := client.List(context.TODO(), req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			logger.Logger.Error("Failed to list VM Instances", "error", err.Error())
			break
		}
		instances = append(instances, *resp)
	}

	return instances
}
