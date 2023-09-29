package signalfx

import (
	"context"
	"github.com/signalfx/signalfx-go/analytics_search"
	"github.com/signalfx/signalfx-go/graphql"
	"time"
)

func (c *Client) StartAnalyticsSearch(ctx context.Context, startTime time.Time, endTime time.Time, traceFilter *analytics_search.TraceFilter, spanFilters []analytics_search.SpanFilter) (string, error) {
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

	query := "query StartAnalyticsSearch($parameters: JSON!) {\n  startAnalyticsSearch(parameters: $parameters)\n}\n"

	graphqlRequest := graphql.Request{
		OperationName: "StartAnalyticsSearch",
		Variables: graphql.Variables{
			Parameters: graphql.Parameters{
				SectionsParameters: []graphql.SectionParams{
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
	err := c.GraphQLRequest(context.Background(), graphqlRequest, &startAnalyticsSearchResponseData)
	if err != nil {
		return "", err
	}

	return startAnalyticsSearchResponseData.Data.StartAnalyticsSearch.JobID, nil
}
