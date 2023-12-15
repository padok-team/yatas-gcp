package iam

import (
	"sync"

	"github.com/padok-team/yatas-gcp/internal"
	"github.com/padok-team/yatas/plugins/commons"
)

type PermissionsBySA struct {
	SA          string
	Permissions []string
}

func (p *PermissionsBySA) GetID() string {
	return p.SA
}

func RunChecks(wa *sync.WaitGroup, account internal.GCPAccount, c *commons.Config, queue chan []commons.Check) {
	var checkConfig commons.CheckConfig
	checkConfig.Init(c)
	var checks []commons.Check

	permissionsBySA := GetPermissionsByServiceAccounts(account)

	iamPermissionsChecks := []commons.CheckDefinition{
		{
			Title:          "GCP_IAM_001",
			Description:    "Service accounts cannot escalate privileges",
			Categories:     []string{"Security", "Good Practice"},
			ConditionFn:    SACannotEscalatePrivileges,
			SuccessMessage: "Service Account cannot escalate privileges. Check the permissions of this service account for extra caution. More on https://rhinosecuritylabs.com/gcp/privilege-escalation-google-cloud-platform-part-1/",
			FailureMessage: "Service Account can escalate privileges! Check the permissions of this service account with caution. More on https://rhinosecuritylabs.com/gcp/privilege-escalation-google-cloud-platform-part-1/",
		},
	}

	var permissionsResources []commons.Resource
	for sa, permissions := range permissionsBySA {
		permissionsResources = append(permissionsResources, &PermissionsBySA{SA: sa, Permissions: permissions})
	}
	commons.AddChecks(&checkConfig, iamPermissionsChecks)
	go commons.CheckResources(checkConfig, permissionsResources, iamPermissionsChecks)

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
