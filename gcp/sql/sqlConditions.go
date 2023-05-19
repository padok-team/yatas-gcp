package sql

import (
	"github.com/padok-team/yatas/plugins/commons"
)

func SQLInstanceIsRegional(resource commons.Resource) bool {
	sqlInstance, ok := resource.(*SQLInstance)
	if !ok {
		return false
	}

	return sqlInstance.Instance.Settings.AvailabilityType == "REGIONAL"
}

func SQLInstanceBackupWithPITREnabled(resource commons.Resource) bool {
	sqlInstance, ok := resource.(*SQLInstance)
	if !ok {
		return false
	}

	if sqlInstance.Instance.Settings != nil && sqlInstance.Instance.Settings.BackupConfiguration != nil {
		// Backup must be enabled and PITR also
		return sqlInstance.Instance.Settings.BackupConfiguration.Enabled &&
			sqlInstance.Instance.Settings.BackupConfiguration.PointInTimeRecoveryEnabled
	} else {
		return false
	}
}
