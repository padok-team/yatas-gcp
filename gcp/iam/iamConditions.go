package iam

import (
	"regexp"
	"slices"
	"time"

	"github.com/padok-team/yatas/plugins/commons"
)

// Source:
// - https://rhinosecuritylabs.com/gcp/privilege-escalation-google-cloud-platform-part-1/
// - https://rhinosecuritylabs.com/cloud-security/privilege-escalation-google-cloud-platform-part-2/
var permissionsThatCanEscalatePrivileges = []string{
	"deploymentmanager.deployments.create",
	"iam.roles.update",
	"iam.serviceAccounts.getAccessToken",
	"iam.serviceAccountKeys.create",
	"iam.serviceAccounts.implicitDelegation",
	"iam.serviceAccounts.signBlob",
	"iam.serviceAccounts.signJwt",
	"iam.serviceAccounts.actAs",
	"orgpolicy.policy.set",
	"storage.hmacKeys.create",
	"serviceusage.apiKeys.create",
	"serviceusage.apiKeys.list",
}

func isAbleToModifyIAMPolicy(permission string) bool {
	// Source: https://rhinosecuritylabs.com/cloud-security/privilege-escalation-google-cloud-platform-part-2/
	pattern := regexp.MustCompile(`^.*\.setIamPolicy$`)
	return pattern.MatchString(permission)
}

func SACannotEscalatePrivileges(resource commons.Resource) bool {
	permissionsBySA, ok := resource.(*PermissionsBySA)
	if !ok {
		return false
	}

	// check if there is a permission that can escalate privileges
	for _, permission := range permissionsBySA.Permissions {
		if isAbleToModifyIAMPolicy(permission) || slices.Contains(permissionsThatCanEscalatePrivileges, permission) {
			return false
		}
	}

	return true
}

func SAKeysNotOlderThan90Days(resource commons.Resource) bool {
	saKey, ok := resource.(*SAKey)
	if !ok {
		return false
	}

	t := time.Unix(saKey.Key.ValidAfterTime.Seconds, 0)
	daysDiff := time.Since(t).Hours() / 24
	return daysDiff <= 90
}
