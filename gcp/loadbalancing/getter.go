package loadbalancing

import (
	"context"
	"net/url"
	"path"
	"strings"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"github.com/padok-team/yatas-gcp/internal"
	"github.com/padok-team/yatas-gcp/logger"
	"google.golang.org/api/iterator"
)

// Get all the forwarding rules of the account, with a target of type TargetHTTPSProxy
func GetForwardingRulesHTTPS(account internal.GCPAccount, region string) []computepb.ForwardingRule {
	if region == "global" {
		return getGlobalForwardingRulesHTTPS(account)
	}

	ctx := context.Background()
	client, err := compute.NewForwardingRulesRESTClient(ctx)
	if err != nil {
		logger.Logger.Error("Failed to create ForwardingRulesClient client", "error", err)
	}
	defer client.Close()

	req := &computepb.ListForwardingRulesRequest{
		Project: account.Project,
		Region:  region,
	}
	var rules []computepb.ForwardingRule
	it := client.List(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			logger.Logger.Error("Failed to list Forwarding rules", "error", err.Error())
			break
		}
		// Keep only the forwarding rules with a target of type TargetHTTPSProxy
		// TODO: Also support the SSLProxy targets
		if resp.Target != nil && strings.Contains(*resp.Target, "targetHttpsProxies") {
			rules = append(rules, *resp)
		}
	}

	return rules
}

// Private function called by GetForwardingRulesHTTPS for getting the global forwarding rules.
// This had to be a separate function because the API expose different clients and interfaces
// for global forwarding rules.
func getGlobalForwardingRulesHTTPS(account internal.GCPAccount) []computepb.ForwardingRule {
	ctx := context.Background()
	client, err := compute.NewGlobalForwardingRulesRESTClient(ctx)
	if err != nil {
		logger.Logger.Error("Failed to create GlobalForwardingRulesRESTClient client", "error", err)
	}
	defer client.Close()

	req := &computepb.ListGlobalForwardingRulesRequest{
		Project: account.Project,
	}
	var rules []computepb.ForwardingRule
	it := client.List(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			logger.Logger.Error("Failed to list Forwarding rules", "error", err.Error())
			break
		}
		// Keep only the forwarding rules with a target of type TargetHTTPSProxy
		if resp.Target != nil && strings.Contains(*resp.Target, "targetHttpsProxies") {
			rules = append(rules, *resp)
		}
	}

	return rules
}

// Extract the resource name from the resource URL (last path part)
func extractResourceName(resourceUrl string) string {
	u, err := url.Parse(resourceUrl)
	if err != nil {
		logger.Logger.Error("Failed to parse resource URL", "url", resourceUrl, "error", err.Error())
		return ""
	}

	return path.Base(u.Path)
}

// Get the SSL certificates associated with a forwarding rule
func GetForwardingRuleSSLCertificate(account internal.GCPAccount, forwardingRule computepb.ForwardingRule) []computepb.SslCertificate {
	// First get the TargetHttpsProxy associated with the forwarding rule
	ctx := context.Background()
	client, err := compute.NewTargetHttpsProxiesRESTClient(ctx)
	if err != nil {
		logger.Logger.Error("Failed to create TargetHttpsProxiesRESTClient client", "error", err)
	}
	defer client.Close()

	req := &computepb.GetTargetHttpsProxyRequest{
		Project:          account.Project,
		TargetHttpsProxy: extractResourceName(*forwardingRule.Target),
	}

	proxy, err := client.Get(ctx, req)
	if err != nil {
		logger.Logger.Error("Failed to get TargetHttpsProxy", "proxy", extractResourceName(*forwardingRule.Target), "error", err.Error())
		return []computepb.SslCertificate{}
	}

	// Second get the SSL certificates associated with the TargetHttpsProxy
	certificates := []computepb.SslCertificate{}

	sslClient, err := compute.NewSslCertificatesRESTClient(ctx)
	if err != nil {
		logger.Logger.Error("Failed to create SslCertificatesRESTClient client", "error", err)
	}
	defer sslClient.Close()

	for _, certUrl := range proxy.SslCertificates {
		sslReq := &computepb.GetSslCertificateRequest{
			Project:        account.Project,
			SslCertificate: extractResourceName(certUrl),
		}
		sslCert, err := sslClient.Get(ctx, sslReq)
		if err != nil {
			logger.Logger.Error("Failed to get SslCertificate", "error", err.Error())
		}
		certificates = append(certificates, *sslCert)
	}

	return certificates
}
