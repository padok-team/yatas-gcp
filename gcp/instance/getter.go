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
	it := client.List(ctx, req)
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
func GetInstances(account internal.GCPAccount, computeZone string) []computepb.Instance {
	ctx := context.Background()
	client, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		logger.Logger.Error("Failed to create Instance client", "error", err)
	}
	defer client.Close()

	req := &computepb.ListInstancesRequest{
		Project: account.Project,
		Zone:    computeZone,
	}
	var instances []computepb.Instance
	it := client.List(ctx, req)
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

func GetDisks(account internal.GCPAccount, computeZone string) []computepb.Disk {
	ctx := context.Background()
	client, err := compute.NewDisksRESTClient(ctx)
	if err != nil {
		logger.Logger.Error("Failed to create Disk client", "error", err)
	}
	defer client.Close()

	req := &computepb.ListDisksRequest{
		Project: account.Project,
		Zone:    computeZone,
	}
	var disks []computepb.Disk
	it := client.List(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			logger.Logger.Error("Failed to list Disks", "error", err.Error())
			break
		}
		disks = append(disks, *resp)
	}

	return disks
}

// Get all the instance groups (of type managed) of the account (zonal and regional)
func GetInstanceGroupsAllZones(account internal.GCPAccount) []computepb.InstanceGroupManager {
	regions := account.ComputeRegions
	zones := GetComputeZones(account)
	ctx := context.Background()

	regionClient, err := compute.NewRegionInstanceGroupManagersRESTClient(ctx)
	if err != nil {
		logger.Logger.Error("Failed to create Regional Instance Group client", "error", err)
	}
	defer regionClient.Close()

	zoneClient, err := compute.NewInstanceGroupManagersRESTClient(ctx)
	if err != nil {
		logger.Logger.Error("Failed to create Zonal Instance Group client", "error", err)
	}
	defer zoneClient.Close()

	var instanceGroups []computepb.InstanceGroupManager

	for _, region := range regions {
		req := &computepb.ListRegionInstanceGroupManagersRequest{
			Project: account.Project,
			Region:  region,
		}
		it := regionClient.List(ctx, req)
		for {
			resp, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				logger.Logger.Error("Failed to list Regional Instance Groups", "error", err.Error())
				break
			}
			instanceGroups = append(instanceGroups, *resp)
		}
	}

	for _, zone := range zones {
		req := &computepb.ListInstanceGroupManagersRequest{
			Project: account.Project,
			Zone:    zone,
		}
		it := zoneClient.List(ctx, req)
		for {
			resp, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				logger.Logger.Error("Failed to list Zonal Instance Groups", "error", err.Error())
				break
			}
			instanceGroups = append(instanceGroups, *resp)
		}
	}

	return instanceGroups
}
