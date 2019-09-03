// package main shows a basic usage pattern of the SiganlFlow client.
package main

import (
	"fmt"
	"log"
	"os"

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

	program := os.Getenv("SIGNALFLOW_PROGRAM")
	if program == "" {
		program = "data('cpu.utilization').publish()"
	}

	comp, err := c.Execute(&signalflow.ExecuteRequest{
		Program: program,
	})
	if err != nil {
		log.Printf("Could not send execute request: %v", err)
		return
	}

	fmt.Printf("Resolution: %v\n", comp.Resolution())
	fmt.Printf("Max Delay: %v\n", comp.MaxDelay())
	fmt.Printf("Detected Lag: %v\n", comp.Lag())

	go func() {
		for msg := range comp.Expirations() {
			fmt.Printf("Got expiration notice for TSID %s", msg.TSID)
		}
	}()

	for msg := range comp.Data() {
		// This will run as long as there is data, or until the websocket gets
		// disconnected.
		if len(msg.Payloads) == 0 {
			fmt.Printf("\rNo data available")
			continue
		}
		for _, pl := range msg.Payloads {
			meta := comp.TSIDMetadata(pl.TSID)
			fmt.Printf("%s %v: %v\n", meta.OriginatingMetric, meta.CustomProperties, pl.Value())
		}
		fmt.Println("")
	}

	err = comp.Err()
	if err != nil {
		log.Printf("Error: %v", comp.Err())
	} else {
		log.Printf("Job completed")
	}
}
