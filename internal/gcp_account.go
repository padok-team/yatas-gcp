package internal

import (
	"google.golang.org/api/storage/v1"
)

type GCP_Account struct {
	Project string `yaml:"project"` // Name of the account in the reports
}

type Client_Account struct {
	Client      *storage.Service
	Gcp_account GCP_Account
}
