package loadbalancing

import (
	"strings"
	"sync"

	"cloud.google.com/go/compute/apiv1/computepb"
	"github.com/padok-team/yatas-gcp/internal"
	"github.com/padok-team/yatas/plugins/commons"
)

type GCSForwardingRulesWithSSL struct {
	ForwardingRule computepb.ForwardingRule
	Certificates   []computepb.SslCertificate
}

func (f *GCSForwardingRulesWithSSL) GetID() string {
	certNames := []string{}
	for _, cert := range f.Certificates {
		certNames = append(certNames, *cert.Name)
	}
	name := *f.ForwardingRule.Name + " (SSL certs: " + strings.Join(certNames, ", ") + ")"
	return name
}

func RunChecks(wa *sync.WaitGroup, account internal.GCPAccount, c *commons.Config, queue chan []commons.Check) {
	var checkConfig commons.CheckConfig
	checkConfig.Init(c)
	var checks []commons.Check

	forwardingRules := []computepb.ForwardingRule{}
	// Get all the forwarding rules (all regions) with a target of type HTTPSProxy
	for _, region := range account.ComputeRegions {
		rules := GetForwardingRulesHTTPS(account, region)
		forwardingRules = append(forwardingRules, rules...)
	}

	// We add the global region because some forwarding rules may use it.
	forwardingRules = append(forwardingRules, GetForwardingRulesHTTPS(account, "global")...)

	// For each forwarding rule, get the SSL certificates attached to it
	forwardingRulesWithSSL := []GCSForwardingRulesWithSSL{}
	for _, forwardingRule := range forwardingRules {
		// Get the SSL certificates attached to the forwarding rule
		certificates := GetForwardingRuleSSLCertificate(account, forwardingRule)
		forwardingRulesWithSSL = append(forwardingRulesWithSSL, GCSForwardingRulesWithSSL{
			ForwardingRule: forwardingRule,
			Certificates:   certificates,
		})
	}

	loadbalancingChecks := []commons.CheckDefinition{
		{
			Title:          "GCP_LB_001",
			Description:    "Check if SSL certificates attached to HTTPS forwarding rules are in auto-renewed (managed mode)",
			Categories:     []string{"Security", "Good Practice"},
			ConditionFn:    SSLCertificatesAreManaged,
			FailureMessage: "Forwarding rule has not auto-renewal for one of its SSL certificates",
			SuccessMessage: "Forwarding rule has auto-renewal for all its SSL certificates",
		},
	}

	resources := make([]commons.Resource, len(forwardingRulesWithSSL))
	for i := range forwardingRulesWithSSL {
		resources[i] = &forwardingRulesWithSSL[i]
	}

	commons.AddChecks(&checkConfig, loadbalancingChecks)
	go commons.CheckResources(checkConfig, resources, loadbalancingChecks)

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
