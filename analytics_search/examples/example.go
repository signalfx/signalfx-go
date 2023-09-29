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
	search, err := client.StartTraceSearch(context.Background(), time.Now().Add(-10*time.Minute), time.Now(), &analytics_search.TraceFilter{
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

	isReady := false
	var traces []analytics_search.LegacyTraceExample
	for !isReady {
		isReady, traces, err = client.GetTraceSearch(context.Background(), search)
		if err != nil {
			panic(err)
		}
		time.Sleep(1 * time.Second)
	}

	fmt.Printf("%+v\n", search)
	fmt.Printf("%+v\n", traces)
}
