package functions

import (
	"cloud.google.com/go/functions/apiv2/functionspb"
	"github.com/padok-team/yatas-gcp/internal"
	"github.com/padok-team/yatas/plugins/commons"
)

func CloudFunctionIsInternal(resource commons.Resource) bool {
	fun, ok := resource.(*CloudFunction)
	if !ok {
		return false
	}
	return fun.Function.ServiceConfig.IngressSettings == functionspb.ServiceConfig_ALLOW_INTERNAL_ONLY ||
		fun.Function.ServiceConfig.IngressSettings == functionspb.ServiceConfig_ALLOW_INTERNAL_AND_GCLB
}

func CloudFunctionIsNotUsingDefaultSA(resource commons.Resource) bool {
	fun, ok := resource.(*CloudFunction)
	if !ok {
		return false
	}
	sa := fun.Function.ServiceConfig.ServiceAccountEmail
	return sa != "" && !internal.IsDefaultComputeEngineSA(sa)
}

func CloudFunctionsDoesNotHaveSecretInEnv(resource commons.Resource) bool {
	fun, ok := resource.(*CloudFunction)
	if !ok {
		return false
	}
	vars := fun.Function.ServiceConfig.EnvironmentVariables
	for name, value := range vars {
		if value != "" && internal.MayBeSensitive(name, value) {
			return false
		}
	}
	return true
}
