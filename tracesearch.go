package signalfx

import (
	"context"
	"errors"
	"github.com/signalfx/signalfx-go/analytics_search"
	"github.com/signalfx/signalfx-go/graphql"
	"time"
)

func (c *Client) StartTraceSearch(ctx context.Context, startTime time.Time, endTime time.Time, traceFilter *analytics_search.TraceFilter, spanFilters []analytics_search.SpanFilter) (string, error) {
	sharedParameters := analytics_search.StartAnalyticsSearchGraphQLSharedParameters{
		TimeRangeMillis: analytics_search.TimeRangeMillis{
			Gte: startTime.UnixMilli(),
			Lte: endTime.UnixMilli(),
		},
		Filters: []analytics_search.Filter{
			{
				TraceFilter: traceFilter,
				SpanFilters: spanFilters,
				FilterType:  "traceFilter",
			},
		},
	}

	query := "query StartTraceSearch($parameters: JSON!) {\n  startAnalyticsSearch(parameters: $parameters)\n}\n"

	graphqlRequest := graphql.Request{
		OperationName: "StartTraceSearch",
		Variables: analytics_search.StartAnalyticsSearchVariables{
			Parameters: analytics_search.Parameters{
				SectionsParameters: []analytics_search.SectionParams{
					{
						SectionType: "traceExamples",
						Limit:       1000,
					},
				},
				SharedParameters: sharedParameters},
		},
		Query: query,
	}

	var startAnalyticsSearchResponseData analytics_search.StartAnalyticsSearchGraphQLResponseData
	err := c.GraphQLRequest(ctx, graphqlRequest, &startAnalyticsSearchResponseData)
	if err != nil {
		return "", err
	}

	return startAnalyticsSearchResponseData.Data.StartAnalyticsSearch.JobID, nil
}

func (c *Client) GetTraceSearch(ctx context.Context, jobId string) (bool, []analytics_search.LegacyTraceExample, error) {
	query := `
	query GetTraceSearch($jobId: ID!) {
		getAnalyticsSearch(jobId: $jobId)
	}`

	graphqlRequest := graphql.Request{
		OperationName: "GetTraceSearch",
		Variables: analytics_search.GetAnalyticsSearchVariables{
			JobID: jobId,
		},
		Query: query,
	}

	var getAnalyticsSearchResponseData analytics_search.GetAnalyticsSearchGraphQLResponseData
	err := c.GraphQLRequest(ctx, graphqlRequest, &getAnalyticsSearchResponseData)
	if err != nil {
		return false, nil, err
	}

	// There should only be one section in the response
	if len(getAnalyticsSearchResponseData.Data.GetAnalyticsSearch.Sections) != 1 {
		return false, nil, errors.New("expected exactly one section in the response for getAnalyticsSearch")
	}

	isSearchComplete := getAnalyticsSearchResponseData.Data.GetAnalyticsSearch.Sections[0].IsComplete
	if !isSearchComplete {
		return false, nil, nil
	}

	return true, getAnalyticsSearchResponseData.Data.GetAnalyticsSearch.Sections[0].LegacyTraceExamples, nil
}
