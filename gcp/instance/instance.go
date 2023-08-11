package instance

import (
	"sync"

	"cloud.google.com/go/compute/apiv1/computepb"
	"github.com/padok-team/yatas-gcp/internal"
	"github.com/padok-team/yatas/plugins/commons"
)

func RunChecks(wa *sync.WaitGroup, account internal.GCPAccount, c *commons.Config, queue chan []commons.Check) {
	var checkConfig commons.CheckConfig
	checkConfig.Init(c)
	var checks []commons.Check

	computeZones := GetComputeZones(account)

	// Get all the instances in all the compute zones specified
	var instances []computepb.Instance
	for _, zone := range computeZones {
		instances = append(instances, GetInstances(account, zone)...)
	}

	// Get all the disks in all the compute zones specified
	var disks []computepb.Disk
	for _, zone := range computeZones {
		disks = append(disks, GetDisks(account, zone)...)
	}

	instanceGroups := GetInstanceGroupsAllZones(account)

	instanceChecks := []commons.CheckDefinition{
		{
			Title:          "GCP_VM_001",
			Description:    "Check if VM instance is not using a public IP address",
			Categories:     []string{"Security", "Good Practice"},
			ConditionFn:    InstanceNoPublicIPAttached,
			SuccessMessage: "VM instance is not using a public IP address",
			FailureMessage: "VM instance is using a public IP address",
		},
	}

	diskChecks := []commons.CheckDefinition{
		{
			Title:          "GCP_VM_002",
			Description:    "Check if VM Disk is encrypted with a customer-managed key",
			Categories:     []string{"Security", "Good Practice"},
			ConditionFn:    DiskIsCustomerEncrypted,
			SuccessMessage: "VM Disk is encrypted with a customer-managed key",
			FailureMessage: "VM Disk is not encrypted with a customer-managed key",
		},
	}

	instanceGroupChecks := []commons.CheckDefinition{
		{
			Title:          "GCP_VM_003",
			Description:    "Check VM instance group is multi-zonal",
			Categories:     []string{"Security", "Good Practice"},
			ConditionFn:    InstanceGroupIsMultiZone,
			SuccessMessage: "VM instance group is multi-zonal",
			FailureMessage: "VM instance group is not multi-zonal",
		},
	}

	var resources []commons.Resource
	for _, instance := range instances {
		resources = append(resources, &VMInstance{Instance: instance})
	}
	var diskResources []commons.Resource
	for _, disk := range disks {
		diskResources = append(diskResources, &VMDisk{Disk: disk})
	}
	var instanceGroupResources []commons.Resource
	for _, instanceGroup := range instanceGroups {
		instanceGroupResources = append(instanceGroupResources, &InstanceGroup{InstanceGroup: instanceGroup})
	}

	commons.AddChecks(&checkConfig, instanceChecks)
	commons.AddChecks(&checkConfig, diskChecks)
	commons.AddChecks(&checkConfig, instanceGroupChecks)
	go commons.CheckResources(checkConfig, resources, instanceChecks)
	go commons.CheckResources(checkConfig, diskResources, diskChecks)
	go commons.CheckResources(checkConfig, instanceGroupResources, instanceGroupChecks)

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
