package sql

import (
	"net"

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

func SQLInstanceEncryptedTrafficEnforced(resource commons.Resource) bool {
	sqlInstance, ok := resource.(*SQLInstance)
	if !ok {
		return false
	}

	return sqlInstance.Instance.Settings != nil && sqlInstance.Instance.Settings.IpConfiguration.RequireSsl
}

func isPublicIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}
	return !parsedIP.IsPrivate()
}

func SQLInstanceNotPublicIP(resource commons.Resource) bool {
	sqlInstance, ok := resource.(*SQLInstance)
	if !ok {
		return false
	}

	for _, ip := range sqlInstance.Instance.IpAddresses {
		if isPublicIP(ip.IpAddress) {
			return false
		}
	}

	return true
}
