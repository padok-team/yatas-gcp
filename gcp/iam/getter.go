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
