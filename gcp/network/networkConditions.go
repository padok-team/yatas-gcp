package network

import "github.com/padok-team/yatas/plugins/commons"

const IAP_CIDR = "35.235.240.0/20"

func FirewallRuleOnlyAllowsIAP(resource commons.Resource) bool {
	rule, ok := resource.(*FirewallRule)
	if !ok {
		return false
	}

	// We want IAP to be the only source range
	if len(rule.Rule.GetSourceRanges()) != 1 {
		return false
	}

	return rule.Rule.GetSourceRanges()[0] != IAP_CIDR
}

func FirewallRuleHasSpecificTargets(resource commons.Resource) bool {
	rule, ok := resource.(*FirewallRule)
	if !ok {
		return false
	}

	return len(rule.Rule.GetTargetTags()) > 0 || len(rule.Rule.GetTargetServiceAccounts()) > 0
}
