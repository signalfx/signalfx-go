// Package writer contains a set of components that accept a single type of
// SignalFx data (e.g. datapoints, trace spans) in a simple manner (e.g. an
// input channel) and then sorts out the complexities of sending that data to
// SignalFx's ingest (or agent/gateway) endpoints. They are intended to be used
// when high volumes of data are expected.  Some of the issues that a writer
// should deal with are:
//
// - Batching of data that should be sent to SignalFx.  It is infeasible to
// send every single data item as a single request but too much batching will
// reduce the timeliness of data into the system.
//
// - Buffering of data items while waiting to be transmitted to SignalFx. The
// buffering could use all available memory to the process or have a limit on
// data waiting to be sent, after which point data is dropped.
//
// - If buffering is limited, then the writer must decide what to do when the
// limit is exceeded.  Given the nature of data in our system, newer data is
// usually more valuable than older data, so the writer should not necessarily
// just drop all new incoming data (although this is relatively simple to
// implement), as that would prioritize old data.
//
// - Sending data concurrently to SignalFx.  At large volumes, sending one
// request at a time to ingest/gateway is probably not going to get enough
// throughput, as usually the network and HTTP RTT is the bottleneck at that
// point.
package writer

import (
	"context"

	"github.com/signalfx/golib/v3/datapoint"
)

// Writer is common to both datapoint and span writers
type Writer interface {
	InternalMetrics(prefix string) []*datapoint.Datapoint
	Start(context.Context)
}
