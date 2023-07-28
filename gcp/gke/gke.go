package gke

import (
	"sync"

	containerpb "cloud.google.com/go/container/apiv1/containerpb"

	"github.com/padok-team/yatas-gcp/internal"
	"github.com/padok-team/yatas/plugins/commons"
)

type GKECluster struct {
	Cluster containerpb.Cluster
}

func (c *GKECluster) GetID() string {
	return c.Cluster.Name
}

func RunChecks(wa *sync.WaitGroup, account internal.GCPAccount, c *commons.Config, queue chan []commons.Check) {
	var checkConfig commons.CheckConfig
	checkConfig.Init(c)
	var checks []commons.Check

	clusters := GetClusters(account)

	gkeChecks := []commons.CheckDefinition{
		{
			Title:          "GCP_GKE_001",
			Description:    "GKE Control Plane is Regional (HA)",
			Categories:     []string{"Security", "Good Practice"},
			ConditionFn:    GKEControlPlaneIsRegional,
			SuccessMessage: "GKE Control Plane is Regional (HA)",
			FailureMessage: "GKE Control Plane is Zonal (not HA)",
		},
		{
			Title:          "GCP_GKE_002",
			Description:    "Workload Identity is enabled",
			Categories:     []string{"Security", "Good Practice"},
			ConditionFn:    GKEIsUsingWorkloadIdentity,
			SuccessMessage: "Workload Identity is enabled",
			FailureMessage: "Workload Identity is not enabled",
		},
		{
			Title:          "GCP_GKE_003",
			Description:    "GKE Control Plane does not have a public endpoint",
			Categories:     []string{"Security", "Good Practice"},
			ConditionFn:    GKEIsNotExposedOnPublicEndpoint,
			SuccessMessage: "Control Plane is not exposed on a public endpoint",
			FailureMessage: "Control Plane is exposed on a public endpoint",
		},
	}

	var resources []commons.Resource
	for _, cluster := range clusters {
		resources = append(resources, &GKECluster{Cluster: *cluster})
	}
	commons.AddChecks(&checkConfig, gkeChecks)
	go commons.CheckResources(checkConfig, resources, gkeChecks)

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
