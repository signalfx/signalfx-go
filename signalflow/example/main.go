package main

import (
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/signalfx/signalfx-go/signalflow"
)

func main() {
	var streamURL = os.Getenv("SIGNALFX_STREAM_URL")
	var accessToken = os.Getenv("SIGNALFX_ACCESS_TOKEN")

	c, err := signalflow.NewClient(
		signalflow.StreamURL(streamURL),
		signalflow.AccessToken(accessToken))
	if err != nil {
		log.Printf("Error creating client: %v", err)
		return
	}

	ch, err := c.Execute(&signalflow.ExecuteRequest{
		Program: "data('cpu.utilization').publish()",
	})
	if err != nil {
		log.Printf("Could not send execute request: %v", err)
		return
	}

	for msg := range ch.Messages() {
		spew.Dump(msg)
		if dm, ok := msg.(*signalflow.DataMessage); ok {
			for _, pl := range dm.Payloads {
				log.Printf("value: %v", pl.Value())
			}
		}
	}
}
