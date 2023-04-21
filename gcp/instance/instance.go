package instance

import (
	"context"
	"fmt"
	"strings"
	"sync"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"github.com/padok-team/yatas-gcp/internal"
	"github.com/padok-team/yatas-gcp/logger"
	"github.com/padok-team/yatas/plugins/commons"
)

type VMInstance struct {
	Instance computepb.Instance
}

func (i *VMInstance) GetID() string {
	zoneURLSplit := strings.Split(i.Instance.GetZone(), "/")
	zoneName := zoneURLSplit[len(zoneURLSplit)-1]
	return fmt.Sprintf("%s/%s (%d)", zoneName, i.Instance.GetName(), i.Instance.GetId())
}

func RunChecks(wa *sync.WaitGroup, account internal.GCPAccount, c *commons.Config, queue chan []commons.Check) {
	var checkConfig commons.CheckConfig
	checkConfig.Init(c)
	var checks []commons.Check

	ctx := context.Background()
	client, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		logger.Logger.Error("Failed to create Instance client", "error", err)
	}
	defer client.Close()

	computeZones := GetComputeZones(account)

	// Get all the instances in all the compute zones specified
	var instances []computepb.Instance
	for _, zone := range computeZones {
		instances = append(instances, GetInstances(account, client, zone)...)
	}

	instanceChecks := []commons.CheckDefinition{
		{
			Title:          "GCP_VM_001",
			Description:    "Check if VM instance is not using a public IP address.",
			Categories:     []string{"Security", "Good Practice"},
			ConditionFn:    InstanceNoPublicIPAttached,
			SuccessMessage: "VM instance is not using a public IP address",
			FailureMessage: "VM instance is using a public IP address",
		},
	}

	var resources []commons.Resource
	for _, instance := range instances {
		resources = append(resources, &VMInstance{Instance: instance})
	}

	commons.AddChecks(&checkConfig, instanceChecks)
	go commons.CheckResources(checkConfig, resources, instanceChecks)

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
