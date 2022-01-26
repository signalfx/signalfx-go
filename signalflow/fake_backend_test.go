package signalflow

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/signalfx/signalfx-go/idtool"
	"github.com/signalfx/signalfx-go/signalflow/messages"
	"github.com/stretchr/testify/assert"
)

const program = "testflow"

type testCase struct {
	timeSeriesProperties []map[string]string
	name                 string
	startMs              int64 // epoch time ms to start the metrics
	stopMs               int64 // epoch time ms to stop the metrics
	resolutionSecs       int64 // seconds gap between the metrics
	expectedTimestamps   []int64
	numberOfSfxClients   int // count of SFX clients to connect to fakebackend
}

func TestFakeBackend(t *testing.T) {

	now := time.Now()
	testCases := []testCase{
		{
			timeSeriesProperties: []map[string]string{},
			name:                 "no metrics with one resolution window",
			stopMs:               now.UnixNano() / (1000 * 1000),
			startMs:              now.Add(-2*time.Second).UnixNano() / (1000 * 1000),
			resolutionSecs:       2,
			expectedTimestamps: []int64{
				now.Add(-2*time.Second).UnixNano() / (1000 * 1000),
				now.UnixNano() / (1000 * 1000),
			},
			numberOfSfxClients: 2,
		},
		{
			timeSeriesProperties: []map[string]string{
				map[string]string{
					"dim1": "val1",
					"dim2": "val2",
				},
				map[string]string{
					"dim1": "val1",
				},
			},
			name:           "some metrics across 2 resolution windows",
			stopMs:         now.UnixNano() / (1000 * 1000),
			startMs:        now.Add(-4*time.Second).UnixNano() / (1000 * 1000),
			resolutionSecs: 2,
			expectedTimestamps: []int64{
				now.Add(-4*time.Second).UnixNano() / (1000 * 1000),
				now.Add(-2*time.Second).UnixNano() / (1000 * 1000),
				now.UnixNano() / (1000 * 1000),
			},
			numberOfSfxClients: 2,
		},
	}

	for _, testCase := range testCases {
		fakeBackend := NewRunningFakeBackend()
		tsids := []idtool.ID{}
		for _, _ = range testCase.timeSeriesProperties {
			tsids = append(tsids, idtool.ID(rand.Int63()))
		}
		for i, ts := range testCase.timeSeriesProperties {
			fakeBackend.AddTSIDMetadata(tsids[i], &messages.MetadataProperties{
				Metric:           program,
				CustomProperties: ts,
			})
			fakeBackend.SetTSIDFloatData(tsids[i], 0)
		}

		fakeBackend.AddProgramTSIDs(program, tsids)

		// connect N clients so we can prove the fakebackend is not killed by the first client disconnecting
		for i := 1; i <= testCase.numberOfSfxClients; i++ {
			sfxClient, _ := NewClient(StreamURL(fakeBackend.URL()), AccessToken(fakeBackend.AccessToken))
			processClient(t, sfxClient, testCase, i)
		}

	}

}

func processClient(t *testing.T, sfxClient *Client, testCase testCase, connectionCount int) {

	data, _ := sfxClient.Execute(&ExecuteRequest{
		Program:      program,
		StartMs:      testCase.startMs,
		StopMs:       testCase.stopMs,
		ResolutionMs: testCase.resolutionSecs * 1000,
	})
	timestamps := []int64{}
	datapointCount := 0
	for msg := range data.Data() {
		timestamps = append(timestamps, int64(msg.TimestampMillis))
		datapoints := []map[string]string{}
		for _, pl := range msg.Payloads {
			meta := data.TSIDMetadata(pl.TSID)
			dims := map[string]string{}
			for k, v := range meta.CustomProperties {
				dims[k] = v
			}
			datapoints = append(datapoints, dims)
			datapointCount++
		}
		// the datapoints should be always the same the fed in mts
		assert.Equal(t, testCase.timeSeriesProperties, datapoints, testCase.name+": datapoints are wrong on connection "+strconv.Itoa(connectionCount))

		if data.IsFinished() {
			sfxClient.Close()
			break
		}
	}

	assert.Equal(t, testCase.expectedTimestamps, timestamps, testCase.name+": timestamps in metrics are wrong on connection "+strconv.Itoa(connectionCount))
	// the number of datapoints should be the number of resolution windows multiplied by the number of MTS in each timestamp payload
	assert.Equal(t, len(testCase.expectedTimestamps)*len(testCase.timeSeriesProperties), datapointCount, testCase.name+": amount of datapoints unexpected on connection "+strconv.Itoa(connectionCount))
	fmt.Println("finisged " + strconv.Itoa(connectionCount))
}
