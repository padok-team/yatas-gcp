package sql

import (
	"context"

	"github.com/padok-team/yatas-gcp/internal"
	"github.com/padok-team/yatas-gcp/logger"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

// Get all the SQL instances of the account
func GetDBInstances(account internal.GCPAccount) []*sqladmin.DatabaseInstance {
	ctx := context.Background()
	client, err := sqladmin.NewService(ctx)
	if err != nil {
		logger.Logger.Error("Failed to create SQL client", "error", err)
	}

	instances, err := client.Instances.List(account.Project).Do()
	if err != nil {
		logger.Logger.Error("Failed to list SQL instances", "error", err)
	}

	return instances.Items
}
