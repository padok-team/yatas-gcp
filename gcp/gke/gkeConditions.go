package gke

import (
	"regexp"

	"github.com/padok-team/yatas/plugins/commons"
)

// A region name is always composed of two parts separated by a dash, and ends with a number.
// For example us-central1.
func isValidRegionName(name string) bool {
	regionPattern := regexp.MustCompile(`^[a-z]+-[a-z]+\d+$`)
	return regionPattern.MatchString(name)
}

func GKEControlPlaneIsRegional(resource commons.Resource) bool {
	cluster, ok := resource.(*GKECluster)
	if !ok {
		return false
	}

	return isValidRegionName(cluster.Cluster.Location)
}

func GKEIsUsingWorkloadIdentity(resource commons.Resource) bool {
	cluster, ok := resource.(*GKECluster)
	if !ok {
		return false
	}

	return cluster.Cluster.WorkloadIdentityConfig != nil && cluster.Cluster.WorkloadIdentityConfig.WorkloadPool != ""
}
