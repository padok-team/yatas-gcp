package sql

import (
	"testing"

	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

type fakeResource struct{}

func (f fakeResource) GetID() string {
	return "fake-resource"
}

func TestSQLInstanceIsRegional(t *testing.T) {
	// Test when resource is not an SQLInstance
	if SQLInstanceIsRegional(&fakeResource{}) {
		t.Error("Expected SQLInstanceIsRegional to return false for non-SQLInstance resource")
	}

	// Test when instance's availability type is REGIONAL
	instance := &SQLInstance{
		Instance: sqladmin.DatabaseInstance{
			Settings: &sqladmin.Settings{
				AvailabilityType: "REGIONAL",
			},
		},
	}
	if !SQLInstanceIsRegional(instance) {
		t.Error("Expected SQLInstanceIsRegional to return true for SQLInstance with availability type REGIONAL")
	}

	// Test when instance's availability type is not REGIONAL
	instance.Instance.Settings.AvailabilityType = "ZONAL"
	if SQLInstanceIsRegional(instance) {
		t.Error("Expected SQLInstanceIsRegional to return false for SQLInstance with availability type ZONAL")
	}
}

func TestSQLInstanceBackupWithPITREnabled(t *testing.T) {
	// Test when resource is not an SQLInstance
	if SQLInstanceBackupWithPITREnabled(&fakeResource{}) {
		t.Error("Expected SQLInstanceBackupWithPITREnabled to return false for non-SQLInstance resource")
	}

	// Test when backup is not enabled
	instance := &SQLInstance{
		Instance: sqladmin.DatabaseInstance{
			Settings: &sqladmin.Settings{
				BackupConfiguration: &sqladmin.BackupConfiguration{
					Enabled: false,
				},
			},
		},
	}
	if SQLInstanceBackupWithPITREnabled(instance) {
		t.Error("Expected SQLInstanceBackupWithPITREnabled to return false when backup is not enabled")
	}

	// Test when PITR is not enabled
	instance.Instance.Settings.BackupConfiguration.Enabled = true
	instance.Instance.Settings.BackupConfiguration.PointInTimeRecoveryEnabled = false
	if SQLInstanceBackupWithPITREnabled(instance) {
		t.Error("Expected SQLInstanceBackupWithPITREnabled to return false when PITR is not enabled")
	}

	// Test when backup and PITR are enabled
	instance.Instance.Settings.BackupConfiguration.PointInTimeRecoveryEnabled = true
	if !SQLInstanceBackupWithPITREnabled(instance) {
		t.Error("Expected SQLInstanceBackupWithPITREnabled to return true when backup and PITR are enabled")
	}
}

func TestSQLInstanceEncryptedTrafficEnforced(t *testing.T) {
	// Test when resource is not an SQLInstance
	if SQLInstanceEncryptedTrafficEnforced(&fakeResource{}) {
		t.Error("Expected SQLInstanceEncryptedTrafficEnforced to return false for non-SQLInstance resource")
	}

	// Test when SSL Require is not referenced by the API
	instance := &SQLInstance{
		Instance: sqladmin.DatabaseInstance{
			Settings: &sqladmin.Settings{
				IpConfiguration: &sqladmin.IpConfiguration{
					Ipv4Enabled: true,
				},
			},
		},
	}
	if SQLInstanceEncryptedTrafficEnforced(instance) {
		t.Error("Expected SQLInstanceEncryptedTrafficEnforced to return false when RequireSsl is not referenced by the API")
	}

	// Test when Require SSL is set to false
	instance.Instance.Settings.IpConfiguration.RequireSsl = false
	if SQLInstanceEncryptedTrafficEnforced(instance) {
		t.Error("Expected SQLInstanceEncryptedTrafficEnforced to return false when RequireSsl is not enabled")
	}

	// Test when Require SSL is set to true
	instance.Instance.Settings.IpConfiguration.RequireSsl = true
	if !SQLInstanceEncryptedTrafficEnforced(instance) {
		t.Error("Expected SQLInstanceEncryptedTrafficEnforced to return true when RequireSsl is enabled")
	}
}

func TestSQLInstanceNotPublicIP(t *testing.T) {
	// Test when resource is not an SQLInstance
	if SQLInstanceNotPublicIP(&fakeResource{}) {
		t.Error("Expected SQLInstanceNotPublicIP to return false for non-SQLInstance resource")
	}

	// Test when there is no IP address
	instance := &SQLInstance{
		Instance: sqladmin.DatabaseInstance{
			IpAddresses: []*sqladmin.IpMapping{},
		},
	}
	if !SQLInstanceNotPublicIP(instance) {
		t.Error("Expected SQLInstanceNotPublicIP to return true when there is no IP address")
	}

	// Test when there is a public IP address
	instance.Instance.IpAddresses = append(instance.Instance.IpAddresses, &sqladmin.IpMapping{
		IpAddress: "1.2.3.4",
		Type:      "PRIMARY",
	})
	if SQLInstanceNotPublicIP(instance) {
		t.Error("Expected SQLInstanceNotPublicIP to return false when there is a public IP address")
	}

	// Test when no IP address is public
	instance.Instance.IpAddresses[0].IpAddress = "192.168.128.1"
	if !SQLInstanceNotPublicIP(instance) {
		t.Error("Expected SQLInstanceNotPublicIP to return true when there is only private IP addresses")
	}
}

func TestSQLInstanceIsEncryptedWithKMS(t *testing.T) {
	// Test when resource is not an SQLInstance
	if SQLInstanceIsEncryptedWithKMS(&fakeResource{}) {
		t.Error("Expected SQLInstanceIsEncryptedWithKMS to return false for non-SQLInstance resource")
	}

	// Test when there is no DiskEncryptionConfiguration
	instance := &SQLInstance{
		Instance: sqladmin.DatabaseInstance{},
	}
	if SQLInstanceIsEncryptedWithKMS(instance) {
		t.Error("Expected SQLInstanceIsEncryptedWithKMS to return false when there is no DiskEncryptionConfiguration")
	}

	// Test when there is a DiskEncryptionConfiguration
	instance.Instance.DiskEncryptionConfiguration = &sqladmin.DiskEncryptionConfiguration{}
	if SQLInstanceIsEncryptedWithKMS(instance) {
		t.Error("Expected SQLInstanceIsEncryptedWithKMS to return false when there is not a KMS key")
	}

	// Test when there is a KMS key
	instance.Instance.DiskEncryptionConfiguration = &sqladmin.DiskEncryptionConfiguration{
		KmsKeyName: "projects/123/locations/global/keyRings/123/cryptoKeys/123",
	}
	if !SQLInstanceIsEncryptedWithKMS(instance) {
		t.Error("Expected SQLInstanceIsEncryptedWithKMS to return true when there is a KMS key")
	}
}
