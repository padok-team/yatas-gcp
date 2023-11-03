package cloudrun

import (
	"cloud.google.com/go/run/apiv2/runpb"
	"github.com/padok-team/yatas-gcp/internal"
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

func CloudRunServiceIsNotUsingDefaultSA(resource commons.Resource) bool {
	svc, ok := resource.(*CloudRunService)
	if !ok {
		return false
	}
	sa := svc.Service.GetTemplate().GetServiceAccount()
	return sa != "" && !internal.IsDefaultComputeEngineSA(sa)
}

func CloudRunServiceDoesNotHaveSecretInEnv(resource commons.Resource) bool {
	svc, ok := resource.(*CloudRunService)
	if !ok {
		return false
	}
	template := svc.Service.GetTemplate()
	if template == nil {
		return false
	}
	containers := template.GetContainers()
	for _, container := range containers {
		envVars := container.GetEnv()
		for _, envVar := range envVars {
			name := envVar.GetName()
			// If GetValue returns an empty string, it means that the env var may be a secret ref, so we do not check it
			value := envVar.GetValue()
			if value != "" && internal.MayBeSensitive(name, value) {
				return false
			}
		}
	}
	return true
}
