package cloudrun

import (
	"cloud.google.com/go/run/apiv2/runpb"
	"github.com/padok-team/yatas/plugins/commons"
)

func CloudRunServiceIsInternal(resource commons.Resource) bool {
	svc, ok := resource.(*CloudRunService)
	if !ok {
		return false
	}

	return svc.Service.Ingress == runpb.IngressTraffic_INGRESS_TRAFFIC_INTERNAL_ONLY ||
		svc.Service.Ingress == runpb.IngressTraffic_INGRESS_TRAFFIC_INTERNAL_LOAD_BALANCER
}
