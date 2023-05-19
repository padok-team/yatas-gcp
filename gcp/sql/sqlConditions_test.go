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