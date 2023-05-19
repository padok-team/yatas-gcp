package loadbalancing

import (
	"github.com/padok-team/yatas/plugins/commons"
)

func SSLCertificatesAreManaged(resource commons.Resource) bool {
	forwardingRule, ok := resource.(*GCSForwardingRulesWithSSL)
	if !ok {
		return false
	}

	for _, certificate := range forwardingRule.Certificates {
		if *certificate.Type != "MANAGED" {
			return false
		}
	}
	return true
}
