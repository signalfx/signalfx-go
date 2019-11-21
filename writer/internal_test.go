package writer

import "github.com/signalfx/golib/v3/datapoint"

func findInternalMetricWithName(writer Writer, name string) int {
	dps := writer.InternalMetrics("")
	for _, dp := range dps {
		if dp.Metric == name {
			return int(dp.Value.(datapoint.IntValue).Int())
		}
	}
	panic("internal metric not found: " + name)
}
