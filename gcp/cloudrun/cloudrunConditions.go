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

func mayBeSensitive(name string, value string) bool {
	privateKeyPattern := regexp.MustCompile(`-----BEGIN (RSA|EC|DSA|GPP|OPENSSH) PRIVATE KEY-----`)
	namePattern := regexp.MustCompile(`(key|secret|password|token|private|credential|auth|certificate|cert|pem|ssl|tls|ssh|rsa|ecdsa|dsa|gpp)(?i)`)

	return namePattern.MatchString(name) || privateKeyPattern.MatchString(value)
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
			if value != "" && mayBeSensitive(name, value) {
				return false
			}
		}
	}
	return true
}
