package sql

import (
	"sync"

	"github.com/padok-team/yatas-gcp/internal"
	"github.com/padok-team/yatas/plugins/commons"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

type SQLInstance struct {
	Instance sqladmin.DatabaseInstance
}

func (s *SQLInstance) GetID() string {
	return s.Instance.Name
}

func RunChecks(wa *sync.WaitGroup, account internal.GCPAccount, c *commons.Config, queue chan []commons.Check) {
	var checkConfig commons.CheckConfig
	checkConfig.Init(c)
	var checks []commons.Check

	// Get all the SQL instances of the account
	dbInstances := GetDBInstances(account)

	sqlChecks := []commons.CheckDefinition{
		{
			Title:          "GCP_SQL_001",
			Description:    "Check if SQL Instances are Regional (HA)",
			Categories:     []string{"Security", "Good Practice"},
			ConditionFn:    SQLInstanceIsRegional,
			SuccessMessage: "SQL instance is Regional (HA)",
			FailureMessage: "SQL instance is Zonal (not HA)",
		},
		{
			Title:          "GCP_SQL_002",
			Description:    "Check if SQL Instances have backups enabled with Point in Time Recovery",
			Categories:     []string{"Security", "Good Practice"},
			ConditionFn:    SQLInstanceBackupWithPITREnabled,
			SuccessMessage: "SQL instance has backups enabled with Point in Time Recovery",
			FailureMessage: "SQL instances does not have backups enabled with Point in Time Recovery",
		},
		{
			Title:          "GCP_SQL_003",
			Description:    "Check if SQL Instances have encrypted traffic enforced",
			Categories:     []string{"Security", "Good Practice"},
			ConditionFn:    SQLInstanceEncryptedTrafficEnforced,
			SuccessMessage: "SQL instance has encrypted traffic enforced (require SSL option)",
			FailureMessage: "SQL instances does not have encrypted traffic enforced (require SSL option)",
		},
		{
			Title:          "GCP_SQL_004",
			Description:    "Check if SQL Instances are not exposed with a public IP",
			Categories:     []string{"Security", "Good Practice"},
			ConditionFn:    SQLInstanceNotPublicIP,
			SuccessMessage: "SQL instance is not exposed with a public IP",
			FailureMessage: "SQL instance is exposed with a public IP",
		},
		{
			Title:          "GCP_SQL_005",
			Description:    "Check if SQL Instances are encrypted at rest with a customer-managed key",
			Categories:     []string{"Security", "Good Practice"},
			ConditionFn:    SQLInstanceIsEncryptedWithKMS,
			SuccessMessage: "SQL instance is encrypted at rest with a customer-managed key",
			FailureMessage: "SQL instance is not encrypted at rest with a customer-managed key",
		},
	}

	var resources []commons.Resource
	for _, instance := range dbInstances {
		resources = append(resources, &SQLInstance{Instance: *instance})
	}
	commons.AddChecks(&checkConfig, sqlChecks)
	go commons.CheckResources(checkConfig, resources, sqlChecks)

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
