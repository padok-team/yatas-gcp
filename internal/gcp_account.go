package internal

type GCPAccount struct {
	Project        string   `yaml:"project"`        // Name of the account in the reports
	Organization   string   `yaml:"organization"`   // Organization id to discover projects from
	Folder         string   `yaml:"folder"`         // Folder id to discover projects from
	ComputeRegions []string `yaml:"computeRegions"` // List of regions to make Compute API calls
}
