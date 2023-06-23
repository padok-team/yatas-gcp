package gke

import (
	"testing"

	containerpb "cloud.google.com/go/container/apiv1/containerpb"
)

type fakeResource struct{}

func (f fakeResource) GetID() string {
	return "fake-resource"
}

func TestIsValidRegionName(t *testing.T) {
	// Test valid region names
	validNames := []string{"us-central1", "europe-west1", "asia-southeast2"}
	for _, name := range validNames {
		if !isValidRegionName(name) {
			t.Errorf("Expected region name '%s' to be valid, but it is considered invalid", name)
		}
	}

	// Test invalid region names
	invalidNames := []string{"us-central1-", "-europe-west1", "asia-southeast2x", "us-central", "europe-west9-a"}
	for _, name := range invalidNames {
		if isValidRegionName(name) {
			t.Errorf("Expected region name '%s' to be invalid, but it is considered valid", name)
		}
	}
}

func TestGKEControlPlaneIsRegional(t *testing.T) {
	// Test when resource is not a GKECluster
	if GKEControlPlaneIsRegional(fakeResource{}) {
		t.Error("Expected GKEControlPlaneIsRegional to return false for non-GKECluster resource")
	}

	// Test when cluster's location is a valid region name
	cluster := &GKECluster{
		Cluster: containerpb.Cluster{
			Location: "us-central1",
		},
	}
	if !GKEControlPlaneIsRegional(cluster) {
		t.Error("Expected GKEControlPlaneIsRegional to return true for GKECluster with regional control plane")
	}

	// Test when cluster's location is not a valid region name
	cluster.Cluster.Location = "us"
	if GKEControlPlaneIsRegional(cluster) {
		t.Error("Expected GKEControlPlaneIsRegional to return false for GKECluster with non-regional control plane")
	}
}

func TestGKEIsUsingWorkloadIdentity(t *testing.T) {
	// Test when resource is not a GKECluster
	if GKEIsUsingWorkloadIdentity(fakeResource{}) {
		t.Error("Expected GKEIsUsingWorkloadIdentity to return false for non-GKECluster resource")
	}

	// Test when cluster is using workload identity
	cluster := &GKECluster{
		Cluster: containerpb.Cluster{
			WorkloadIdentityConfig: &containerpb.WorkloadIdentityConfig{
				WorkloadPool: "my-workload-pool",
			},
		},
	}
	if !GKEIsUsingWorkloadIdentity(cluster) {
		t.Error("Expected GKEIsUsingWorkloadIdentity to return true for GKECluster using workload identity")
	}

	// Test when cluster is not using workload identity
	cluster.Cluster.WorkloadIdentityConfig = nil
	if GKEIsUsingWorkloadIdentity(cluster) {
		t.Error("Expected GKEIsUsingWorkloadIdentity to return false for GKECluster not using workload identity")
	}
}

func TestGKEIsNotExposedOnPublicEndpoint(t *testing.T) {
	// Test when resource is not a GKECluster
	if GKEIsNotExposedOnPublicEndpoint(fakeResource{}) {
		t.Error("Expected GKEIsNotExposedOnPublicEndpoint to return false for non-GKECluster resource")
	}

	// Test when private endpoint is enabled and cluster endpoint is not a public IP
	cluster := &GKECluster{
		Cluster: containerpb.Cluster{
			PrivateClusterConfig: &containerpb.PrivateClusterConfig{
				EnablePrivateEndpoint: true,
			},
			Endpoint: "10.0.0.1",
		},
	}
	if !GKEIsNotExposedOnPublicEndpoint(cluster) {
		t.Error("Expected GKEIsNotExposedOnPublicEndpoint to return true for GKECluster with private endpoint and non-public IP endpoint")
	}

	// Test when private endpoint is enabled and cluster endpoint is a public IP
	cluster.Cluster.Endpoint = "35.123.45.67"
	if GKEIsNotExposedOnPublicEndpoint(cluster) {
		t.Error("Expected GKEIsNotExposedOnPublicEndpoint to return false for GKECluster with private endpoint and public IP endpoint")
	}

	// Test when private endpoint is not enabled
	cluster.Cluster.PrivateClusterConfig.EnablePrivateEndpoint = false
	if GKEIsNotExposedOnPublicEndpoint(cluster) {
		t.Error("Expected GKEIsNotExposedOnPublicEndpoint to return false for GKECluster with private endpoint disabled")
	}
}
