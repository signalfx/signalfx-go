/*
 * Navigators API
 */

package navigator

import "github.com/signalfx/signalfx-go/dashboard_group"

type Navigator struct {
	Id                         string                             `json:"id,omitempty"`
	NavigatorCode              string                             `json:"navigatorCode,omitempty"`
	DisplayName                string                             `json:"displayName,omitempty"`
	IdDisplayName              string                             `json:"idDisplayName,omitempty"`
	Created                    int64                              `json:"created,omitempty"`
	Creator                    string                             `json:"creator,omitempty"`
	LastUpdated                int64                              `json:"lastUpdated,omitempty"`
	LastUpdatedBy              string                             `json:"lastUpdatedBy,omitempty"`
	OwnerScope                 string                             `json:"ownerScope,omitempty"`
	PropertyIdentifierTemplate string                             `json:"propertyIdentifierTemplate,omitempty"`
	EntityMetrics              []*Metric                          `json:"entityMetrics,omitempty"`
	InstanceLabel              string                             `json:"instanceLabel,omitempty"`
	SystemTypes                []string                           `json:"systemTypes,omitempty"`
	ImportQualifiers           []*dashboard_group.ImportQualifier `json:"importQualifiers,omitempty"`
	Category                   *Category                          `json:"category,omitempty"`
	DefaultGroupBy             string                             `json:"defaultGroupBy,omitempty"`
	AlertQuery                 string                             `json:"alertQuery,omitempty"`
	ListColumns                []*ListColumn                      `json:"listColumns,omitempty"`
	SummaryMetricLabel         string                             `json:"summaryMetricLabel,omitempty"`
	SummaryMetricProgramText   string                             `json:"summaryMetricProgramText,omitempty"`
	TooltipKeyList             []*TooltipKey                      `json:"tooltipKeyList,omitempty"`
	DashboardDiscoveryQuery    []string                           `json:"dashboardDiscoveryQuery,omitempty"`
	DashboardMtsQuery          string                             `json:"dashboardMtsQuery,omitempty"`
	RequiredProperties         []string                           `json:"requiredProperties,omitempty"`
	AggregateDashboardName     string                             `json:"aggregateDashboardName,omitempty"`
	InstanceDashboardName      string                             `json:"instanceDashboardName,omitempty"`
	DashboardNameMatch         string                             `json:"dashboardNameMatch,omitempty"`
	AggregateDashboards        []string                           `json:"aggregateDashboards,omitempty"`
	InstanceDashboards         []string                           `json:"instanceDashboards,omitempty"`
	InstanceDisplayText        string                             `json:"instanceDisplayText,omitempty"`
}
