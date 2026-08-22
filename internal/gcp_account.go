package internal

type GCPAccount struct {
	Project        string   `yaml:"project"`        // Name of the account in the reports
	ComputeRegions []string `yaml:"computeRegions"` // List of regions to make Compute API calls
	Organization   string   `yaml:"organization"`   // Organization ID to scan projects
	Folder         string   `yaml:"folder"`         // Folder ID to scan projects
}
