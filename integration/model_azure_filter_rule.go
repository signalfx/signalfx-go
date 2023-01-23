/*
 * Integrations API
 */

package integration

// AzureFilterRule defines Azure resource filter rule; see AzureIntegration.ResourceFilterRules for details
type AzureFilterRule struct {
	Filter AzureFilterExpression `json:"filter"`
}

type AzureFilterExpression struct {
	Source string `json:"source"`
}
