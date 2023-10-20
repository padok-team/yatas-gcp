package functions

import (
	"sync"

	"cloud.google.com/go/functions/apiv2/functionspb"
	"github.com/padok-team/yatas-gcp/internal"
	"github.com/padok-team/yatas/plugins/commons"
)

type CloudFunction struct {
	Function *functionspb.Function
}

func (c *CloudFunction) GetID() string {
	return c.Function.Name
}

func RunChecks(wa *sync.WaitGroup, account internal.GCPAccount, c *commons.Config, queue chan []commons.Check) {
	var checkConfig commons.CheckConfig
	checkConfig.Init(c)
	var checks []commons.Check

	functions := GetCloudFunctions(account)

	functionsChecks := []commons.CheckDefinition{
		{
			Title:          "GCP_FUN_001",
			Description:    "CloudFunctions are not directly exposed on the internet",
			Categories:     []string{"Security", "Good Practice"},
			ConditionFn:    CloudFunctionIsInternal,
			SuccessMessage: "Function is exposed internally or through a load balancer",
			FailureMessage: "Function is directly exposed on the internet",
		},
		{
			Title:          "GCP_FUN_002",
			Description:    "CloudFunctions do not use the default Compute Engine service account",
			Categories:     []string{"Security", "Good Practice"},
			ConditionFn:    CloudFunctionIsNotUsingDefaultSA,
			SuccessMessage: "Function is not using the default Compute Engine service account",
			FailureMessage: "Function is using the default Compute Engine service account",
		},
		{
			Title:          "GCP_FUN_003",
			Description:    "CloudFunctions do not have plain text secrets in environment variables",
			Categories:     []string{"Security", "Good Practice"},
			ConditionFn:    CloudFunctionsDoesNotHaveSecretInEnv,
			SuccessMessage: "Function SEEMS to not have plain text secrets in environment variables, check manually",
			FailureMessage: "Function MIGHT have plain text secrets in environment variables, check manually",
		},
	}

	var resources []commons.Resource
	for _, f := range functions {
		resources = append(resources, &CloudFunction{Function: f})
	}
	commons.AddChecks(&checkConfig, functionsChecks)
	go commons.CheckResources(checkConfig, resources, functionsChecks)

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
