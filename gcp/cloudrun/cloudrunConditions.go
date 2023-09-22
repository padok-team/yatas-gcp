package cloudrun

import (
	"regexp"

	"cloud.google.com/go/run/apiv2/runpb"
	"github.com/padok-team/yatas/plugins/commons"
)

// The default Compute Engine service account is always PROJECT_NUMBER-compute@developer.gserviceaccount.com
func isDefaultComputeEngineSA(name string) bool {
	regionPattern := regexp.MustCompile(`^[0-9]+-compute@developer\.gserviceaccount\.com$`)
	return regionPattern.MatchString(name)
}

func CloudRunServiceIsInternal(resource commons.Resource) bool {
	svc, ok := resource.(*CloudRunService)
	if !ok {
		return false
	}
	return svc.Service.Ingress == runpb.IngressTraffic_INGRESS_TRAFFIC_INTERNAL_ONLY ||
		svc.Service.Ingress == runpb.IngressTraffic_INGRESS_TRAFFIC_INTERNAL_LOAD_BALANCER
}

func CloudRunServiceIsNotUsingDefaultSA(resource commons.Resource) bool {
	svc, ok := resource.(*CloudRunService)
	if !ok {
		return false
	}
	sa := svc.Service.GetTemplate().GetServiceAccount()
	return sa != "" && !isDefaultComputeEngineSA(sa)
}
