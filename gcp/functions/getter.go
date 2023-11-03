package functions

import (
	"context"

	functions "cloud.google.com/go/functions/apiv2"
	"cloud.google.com/go/functions/apiv2/functionspb"
	"cloud.google.com/go/iam/apiv1/iampb"
	run "cloud.google.com/go/run/apiv2"
	"github.com/padok-team/yatas-gcp/internal"
	"github.com/padok-team/yatas-gcp/logger"
	"google.golang.org/api/iterator"
)

func GetCloudFunctions(account internal.GCPAccount) []*functionspb.Function {
	// Set up authentication and create the Cloud Functions client.
	ctx := context.Background()
	client, err := functions.NewFunctionClient(ctx)
	if err != nil {
		logger.Logger.Error("Failed to create CloudFunctions client", "error", err)
		return nil
	}
	defer client.Close()

	var functions []*functionspb.Function

	for _, region := range account.ComputeRegions {
		req := &functionspb.ListFunctionsRequest{
			Parent: "projects/" + account.Project + "/locations/" + region,
		}
		it := client.ListFunctions(ctx, req)
		for {
			resp, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				logger.Logger.Error("Failed to list CloudFunctions", "error", err.Error())
				break
			}
			functions = append(functions, resp)
		}
	}

	return functions
}

func getPolicyWithCloudRun(ctx context.Context, resource string) *iampb.Policy {
	c, err := run.NewServicesClient(ctx)
	if err != nil {
		logger.Logger.Error("Failed to create CloudRun Services client", "error", err)
	}
	defer c.Close()

	policy, err := c.GetIamPolicy(ctx, &iampb.GetIamPolicyRequest{
		Resource: resource,
		Options: &iampb.GetPolicyOptions{
			RequestedPolicyVersion: 1,
		},
	})

	if err != nil {
		logger.Logger.Error("Failed to get IAM policy", "resource", resource, "error", err.Error())
		return nil
	}
	return policy
}

func GetCloudFunctionPolicy(function *functionspb.Function) *iampb.Policy {
	ctx := context.Background()
	// If Gen2 use Cloud Run API
	if function.Environment == functionspb.Environment_GEN_2 {
		return getPolicyWithCloudRun(ctx, function.ServiceConfig.Service)
	}

	client, err := functions.NewFunctionRESTClient(ctx)
	if err != nil {
		logger.Logger.Error("Failed to create CloudFunctions Services client", "error", err)
		return nil
	}
	defer client.Close()

	policy, err := client.GetIamPolicy(ctx, &iampb.GetIamPolicyRequest{
		Resource: function.Name,
		Options: &iampb.GetPolicyOptions{
			RequestedPolicyVersion: 1,
		},
	})

	if err != nil {
		logger.Logger.Error("Failed to get IAM policy", "resource", function.Name, "error", err.Error())
		return nil
	}

	return policy
}
