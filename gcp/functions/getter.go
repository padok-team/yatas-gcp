package functions

import (
	"context"

	functions "cloud.google.com/go/functions/apiv2"
	"cloud.google.com/go/functions/apiv2/functionspb"
	"github.com/padok-team/yatas-gcp/internal"
	"github.com/padok-team/yatas-gcp/logger"
	"google.golang.org/api/iterator"
)

func GetCloudFunctions(account internal.GCPAccount) []*functionspb.Function {
	// Set up authentication and create the Cloud Functions client.
	ctx := context.Background()
	client, err := functions.NewFunctionClient(ctx)
	if err != nil {
		logger.Logger.Error("Failed to create CloudRun Services client", "error", err)
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
