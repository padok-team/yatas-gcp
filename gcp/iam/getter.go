package iam

import (
	"context"
	"strings"

	admin "cloud.google.com/go/iam/admin/apiv1"
	"cloud.google.com/go/iam/admin/apiv1/adminpb"
	"cloud.google.com/go/iam/apiv1/iampb"
	resourcemanager "cloud.google.com/go/resourcemanager/apiv3"
	"github.com/padok-team/yatas-gcp/internal"
	"github.com/padok-team/yatas-gcp/logger"
	"google.golang.org/api/iterator"
)

// Return the permissions for the given role
func getRolePermissions(role string, ctx context.Context, c *admin.IamClient) []string {
	req := &adminpb.GetRoleRequest{
		Name: role,
	}
	resp, err := c.GetRole(ctx, req)
	if err != nil {
		logger.Logger.Error("Failed to get permissions for role", "role", role, "error", err)
		return nil
	}

	var permissions []string

	permissions = append(permissions, resp.GetIncludedPermissions()...)

	return permissions
}

// Get a map of service account -> permission(s) for the current project
func GetPermissionsByServiceAccounts(account internal.GCPAccount) map[string][]string {
	ctx := context.Background()

	c, err := resourcemanager.NewProjectsClient(ctx)
	if err != nil {
		logger.Logger.Error("Failed to create projects client", "error", err)
		return nil
	}
	defer c.Close()

	iamClient, err := admin.NewIamClient(ctx)
	if err != nil {
		logger.Logger.Error("Failed to create IAM client", "error", err)
		return nil
	}
	defer iamClient.Close()

	req := &iampb.GetIamPolicyRequest{
		Resource: "projects/" + account.Project,
	}
	resp, err := c.GetIamPolicy(ctx, req)
	if err != nil {
		logger.Logger.Error("Failed to get IAM policy for project", "project", account.Project, "error", err)
		return nil
	}

	// Iterate over the bindings, keep only the ones that are for service accounts, and create a map of service account -> permission(s)
	permissionsBySA := make(map[string][]string)
	for _, binding := range resp.GetBindings() {
		for _, member := range binding.GetMembers() {
			if strings.HasPrefix(member, "serviceAccount:") {
				permissions := getRolePermissions(binding.GetRole(), ctx, iamClient)
				permissionsBySA[member] = append(permissionsBySA[member], permissions...)
			}
		}
	}

	return permissionsBySA
}

// Get the keys associated to a service account
func getKeysOfServiceAccount(sa *adminpb.ServiceAccount, ctx context.Context, c *admin.IamClient) []*adminpb.ServiceAccountKey {
	req := &adminpb.ListServiceAccountKeysRequest{
		Name: sa.GetName(),
	}
	resp, err := c.ListServiceAccountKeys(ctx, req)
	if err != nil {
		logger.Logger.Error("Failed to list service accounts keys", "serviceAccount", sa.GetName(), "error", err)
		return nil
	}

	return resp.GetKeys()
}

// Produce a list of service account keys existing in the GCP project
func GetServiceAccountKeys(account internal.GCPAccount) []*adminpb.ServiceAccountKey {
	ctx := context.Background()

	c, err := admin.NewIamClient(ctx)
	if err != nil {
		logger.Logger.Error("Failed to create IAM client", "error", err)
		return nil
	}
	defer c.Close()

	req := &adminpb.ListServiceAccountsRequest{
		Name: "projects/" + account.Project,
	}
	var saList []*adminpb.ServiceAccount
	it := c.ListServiceAccounts(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			logger.Logger.Error("Failed to list service accounts for project", "project", account.Project, "error", err)
			return nil // TODO
		}

		saList = append(saList, resp)
	}

	// Get all service accounts keys and merge them into a single list
	var keys []*adminpb.ServiceAccountKey
	for _, sa := range saList {
		saKeys := getKeysOfServiceAccount(sa, ctx, c)
		keys = append(keys, saKeys...)
	}

	return keys
}
