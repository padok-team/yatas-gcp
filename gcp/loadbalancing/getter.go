package loadbalancing

import (
	"context"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"github.com/padok-team/yatas-gcp/internal"
	"github.com/padok-team/yatas-gcp/logger"
	"google.golang.org/api/iterator"
)

// Get all the forwarding rules of the account
func GetForwardingRules(account internal.GCPAccount) []computepb.ForwardingRule {
	ctx := context.Background()
	client, err := compute.NewForwardingRulesRESTClient(ctx)
	if err != nil {
		logger.Logger.Error("Failed to create ForwardingRulesClient client", "error", err)
	}
	defer client.Close()

	req := &computepb.ListForwardingRulesRequest{
		Project: account.Project,
	}
	var rules []computepb.ForwardingRule
	it := client.List(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			logger.Logger.Error("Failed to list Forwarding rules", "error", err.Error())
			break
		}
		rules = append(rules, *resp)
	}

	return rules
}
