#!/bin/sh

if ! [ -x "$(command -v genny)" ]; then
	go install github.com/mauricelam/genny
fi

# Generate writer/buffer for datapoints
genny -ast -pkg=writer -in=ring.go -out=../datapoint_buffer.gen.go -imp "github.com/signalfx/golib/v3/datapoint" gen "Instance=Datapoint:datapoint.Datapoint"

genny -ast -pkg=writer -in=writer.go -out=../datapoint_writer.gen.go -imp "github.com/signalfx/golib/v3/datapoint" gen "Instance=Datapoint:datapoint.Datapoint"

sed -i'' -e 's/instance\(s\?\)_/datapoint\1_/' ../datapoint_writer.gen.go


# Generate writer/buffer for spans
genny -ast -pkg=writer -in=ring.go -out=../span_buffer.gen.go -imp "github.com/signalfx/golib/v3/trace" gen "Instance=Span:trace.Span"

genny -ast -pkg=writer -in=writer.go -out=../span_writer.gen.go -imp "github.com/signalfx/golib/v3/trace" gen "Instance=Span:trace.Span"

sed -i'' -e 's/instance\(s\?\)_/trace_span\1_/' ../span_writer.gen.go


# Common cleanup
goimports -w ../*.gen.go
