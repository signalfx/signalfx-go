package main

import (
	"context"
	"fmt"
	"github.com/signalfx/signalfx-go"
	"github.com/signalfx/signalfx-go/analytics_search"
	"net/http"
	"os"
	"time"
)

func main() {
	httpClient := &http.Client{}
	token := os.Getenv("SIGNALFX_API_TOKEN")
	client, err := signalfx.NewClient(token, signalfx.HTTPClient(httpClient),
		signalfx.APIUrl("https://api.eu0.signalfx.com"),
		signalfx.FrontendUrl("https://app.eu0.signalfx.com"),
	)
	if err != nil {
		panic(err)
	}

	// Then do things!
	search, err := client.StartAnalyticsSearch(context.Background(), time.Now().Add(-10*time.Minute), time.Now(), &analytics_search.TraceFilter{
		Tags: []analytics_search.Tag{
			{
				Tag:       "sf_error",
				Operation: "IN",
				Values:    []string{"true"},
			},
		}}, []analytics_search.SpanFilter{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", search)
}
