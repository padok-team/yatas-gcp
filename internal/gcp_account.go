package internal

type GCPAccount struct {
	Project        string   `yaml:"project"`        // Name of the account in the reports
	ComputeRegions []string `yaml:"computeRegions"` // List of regions to make Compute API calls
}
