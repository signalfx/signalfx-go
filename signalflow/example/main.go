// package main shows a basic usage pattern of the SiganlFlow client.
package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/signalfx/signalfx-go/signalflow/v2"
)

func main() {
	var (
		realm       string
		accessToken string
		program     string
		duration    time.Duration
	)

	flag.StringVar(&realm, "realm", "", "SignalFx Realm")
	flag.StringVar(&accessToken, "access-token", "", "SignalFx Org Access Token")
	flag.StringVar(&program, "program", "data('cpu.utilization').count().publish()", "The SignalFlow program to execute")
	flag.DurationVar(&duration, "duration", 30*time.Second, "How long to run the job before sending Stop message")
	flag.Parse()

	if realm == "" || accessToken == "" {
		flag.Usage()
		os.Exit(1)
	}

	c, err := signalflow.NewClient(
		signalflow.StreamURLForRealm(realm),
		signalflow.AccessToken(accessToken),
		signalflow.OnError(func(err error) {
			log.Printf("Error in SignalFlow client: %v", err)
		}))
	if err != nil {
		log.Printf("Error creating client: %v", err)
		return
	}

	log.Printf("Executing program for %v: %s", duration, program)
	comp, err := c.Execute(context.Background(), &signalflow.ExecuteRequest{
		Program: program,
	})
	if err != nil {
		log.Printf("Could not send execute request: %v", err)
		return
	}

	// If you want to limit how long to wait for the resolution metadata to come in you can use a
	// timed context.
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	resolution, err := comp.Resolution(ctx)
	cancel()

	maxDelay, _ := comp.MaxDelay(context.Background())
	lag, _ := comp.Lag(context.Background())

	log.Printf("Resolution: %v (err: %v)\n", resolution, err)
	log.Printf("Max Delay: %v\n", maxDelay)
	log.Printf("Detected Lag: %v\n", lag)

	go func() {
		for msg := range comp.Expirations() {
			log.Printf("Got expiration notice for TSID %s", msg.TSID)
		}
	}()

	go func() {
		time.Sleep(duration)
		if err := comp.Stop(context.Background()); err != nil {
			log.Printf("Failed to stop computation: %v", err)
		}
	}()

	for msg := range comp.Data() {
		// This will run as long as there is data, or until the websocket gets
		// disconnected.
		if len(msg.Payloads) == 0 {
			log.Printf("\rNo data available")
			continue
		}
		for _, pl := range msg.Payloads {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			meta, err := comp.TSIDMetadata(ctx, pl.TSID)
			cancel()
			if err != nil {
				log.Printf("Failed to get metadata for tsid %s: %v", pl.TSID, err)
				continue
			}
			log.Printf("%s (%s) %v %v: %v\n", meta.OriginatingMetric, meta.Metric, meta.CustomProperties, meta.InternalProperties, pl.Value())
		}
		log.Println("")
	}

	err = comp.Err()
	if err != nil {
		log.Printf("Error: %v", comp.Err())
	} else {
		log.Printf("Job completed")
	}
}
