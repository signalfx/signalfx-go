package main

import (
	"context"
	"fmt"
	"github.com/signalfx/signalfx-go"
	"net/http"
	"os"
	"time"
)

func main() {
	httpClient := &http.Client{}
	token := os.Getenv("SIGNALFX_API_TOKEN")
	client, err := signalfx.NewClient(token, signalfx.HTTPClient(httpClient), signalfx.APIUrl("https://api.eu0.signalfx.com"))
	if err != nil {
		panic(err)
	}

	// Then do things!
	events, err := client.GetEvents(context.Background(), "*", time.Now().Add(-15*time.Minute), time.Now(), 100, 0)
	if err != nil {
		panic(err)
	}
	for _, event := range events {
		fmt.Printf("%+v\n", event)
	}
}
