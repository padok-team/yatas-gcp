package network

import (
	"sync"

	"cloud.google.com/go/compute/apiv1/computepb"

	"github.com/padok-team/yatas-gcp/internal"
	"github.com/padok-team/yatas/plugins/commons"
)

type FirewallRule struct {
	Rule *computepb.Firewall
}

func (r *FirewallRule) GetID() string {
	return r.Rule.GetName()
}

func RunChecks(wa *sync.WaitGroup, account internal.GCPAccount, c *commons.Config, queue chan []commons.Check) {
	var checkConfig commons.CheckConfig
	checkConfig.Init(c)
	var checks []commons.Check

	firewallRules := GetFirewallRulesWithPort22(account)

	firewallChecks := []commons.CheckDefinition{
		{
			Title:          "GCP_NET_001",
			Description:    "SSH ingress firewall rules only allow IAP",
			Categories:     []string{"Security", "Good Practice"},
			ConditionFn:    FirewallRuleOnlyAllowsIAP,
			SuccessMessage: "Ingress rule is only allowing IAP",
			FailureMessage: "Ingress rule allow other CIDRs than IAP",
		},
		{
			Title:          "GCP_NET_002",
			Description:    "SSH ingress firewall rules apply on specific tags or service accounts",
			Categories:     []string{"Security", "Good Practice"},
			ConditionFn:    FirewallRuleHasSpecificTargets,
			SuccessMessage: "Ingress rule is targetting tags or service accounts",
			FailureMessage: "Ingress rule is targetting all instances",
		},
	}

	var resources []commons.Resource
	for _, rule := range firewallRules {
		resources = append(resources, &FirewallRule{Rule: rule})
	}
	commons.AddChecks(&checkConfig, firewallChecks)
	go commons.CheckResources(checkConfig, resources, firewallChecks)

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
