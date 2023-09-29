package main

import (
	"context"
	"fmt"
	"github.com/signalfx/signalfx-go"
	"net/http"
	"os"
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
	trace, err := client.GetTrace(context.Background(), "63d16b254ec5e80d")
	if err != nil {
		panic(err)
	}
	for _, span := range trace {
		fmt.Printf("%+v\n", span)
	}
}
