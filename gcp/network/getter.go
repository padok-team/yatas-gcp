package network

import (
	"context"
	"slices"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"github.com/padok-team/yatas-gcp/internal"
	"github.com/padok-team/yatas-gcp/logger"
	"google.golang.org/api/iterator"
)

func keepFirewallRuleSSH(rule *computepb.Firewall) bool {
	if rule.GetDirection() != "INGRESS" {
		return false
	}

	// Iterate over allow and check if port 22 on tcp is open
	for _, allow := range rule.GetAllowed() {
		if allow.GetIPProtocol() == "tcp" && slices.Contains(allow.GetPorts(), "22") {
			return true
		}
	}

	return false
}

// Return all firewall rules that allow ingress with on port TCP/22
func GetFirewallRulesWithPort22(account internal.GCPAccount) []*computepb.Firewall {
	ctx := context.Background()
	client, err := compute.NewFirewallsRESTClient(ctx)
	if err != nil {
		logger.Logger.Error("Failed to create Firewall client", "error", err)
	}
	defer client.Close()

	var rules []*computepb.Firewall

	req := &computepb.ListFirewallsRequest{
		Project: account.Project,
	}
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
		if keepFirewallRuleSSH(resp) {
			rules = append(rules, resp)
		}
	}

	return rules

}
