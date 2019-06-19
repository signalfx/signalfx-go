package signalflow

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSerializeExecuteRequest(t *testing.T) {
	er := ExecuteRequest{
		Program:    "data(cpu.utilization).publish()",
		Start:      time.Unix(5000, 0),
		Stop:       time.Unix(6000, 0),
		Resolution: 5 * time.Second,
		MaxDelay:   3 * time.Second,
		StopMs:     6500,
	}

	serialized, err := json.Marshal(er)
	require.Nil(t, err)
	require.Equal(t, `{"type":"execute","program":"data(cpu.utilization).publish()","channel":"","start":5000000,"stop":6000000,"resolution":5000,"maxDelay":3000,"immediate":false,"timezone":""}`, string(serialized))

}
