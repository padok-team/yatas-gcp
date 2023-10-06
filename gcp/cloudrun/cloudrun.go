package cloudrun

import (
	"sync"

	"cloud.google.com/go/run/apiv2/runpb"
	"github.com/padok-team/yatas-gcp/internal"
	"github.com/padok-team/yatas/plugins/commons"
)

type CloudRunService struct {
	Service runpb.Service
}

func (c *CloudRunService) GetID() string {
	return c.Service.Name
}

func RunChecks(wa *sync.WaitGroup, account internal.GCPAccount, c *commons.Config, queue chan []commons.Check) {
	var checkConfig commons.CheckConfig
	checkConfig.Init(c)
	var checks []commons.Check

	services := GetCloudRunServices(account)

	cloudrunChecks := []commons.CheckDefinition{
		{
			Title:          "GCP_RUN_001",
			Description:    "CloudRun services are not directly exposed on the internet",
			Categories:     []string{"Security", "Good Practice"},
			ConditionFn:    CloudRunServiceIsInternal,
			SuccessMessage: "Service is exposed internally or through a load balancer",
			FailureMessage: "Service is directly exposed on the internet",
		},
		{
			Title:          "GCP_RUN_002",
			Description:    "CloudRun services do not use the default Compute Engine service account",
			Categories:     []string{"Security", "Good Practice"},
			ConditionFn:    CloudRunServiceIsNotUsingDefaultSA,
			SuccessMessage: "Service is not using the default Compute Engine service account",
			FailureMessage: "Service is using the default Compute Engine service account",
		},
		{
			Title:          "GCP_RUN_003",
			Description:    "CloudRun services do not have plain text secrets in environment variables",
			Categories:     []string{"Security", "Good Practice"},
			ConditionFn:    CloudRunServiceDoesNotHaveSecretInEnv,
			SuccessMessage: "Service SEEMS to not have plain text secrets in environment variables, check manually",
			FailureMessage: "Service MIGHT have plain text secrets in environment variables, check manually",
		},
	}

	var resources []commons.Resource
	for _, svc := range services {
		resources = append(resources, &CloudRunService{Service: svc})
	}
	commons.AddChecks(&checkConfig, cloudrunChecks)
	go commons.CheckResources(checkConfig, resources, cloudrunChecks)

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
