package functions

import (
	"cloud.google.com/go/functions/apiv2/functionspb"
	"cloud.google.com/go/iam/apiv1/iampb"
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
	return sa != "" && !internal.IsDefaultComputeEngineSA(sa) && !internal.IsDefaultAppEngineSA(sa, fun.Project)
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

func isPolicyAllowingAllUsers(policy *iampb.Policy) bool {
	for _, binding := range policy.Bindings {
		if binding.Role == "roles/cloudfunctions.invoker" || binding.Role == "roles/run.invoker" {
			for _, member := range binding.Members {
				if member == "allUsers" {
					return true
				}
			}
		}
	}
	return false
}

func CloudFunctionRequireIAMAuthentication(resource commons.Resource) bool {
	fun, ok := resource.(*CloudFunction)
	if !ok {
		return false
	}
	return fun.Policy != nil && !isPolicyAllowingAllUsers(fun.Policy)
}
