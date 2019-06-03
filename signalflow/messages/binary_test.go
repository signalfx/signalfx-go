package messages

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/signalfx/signalfx-go/idtool"
	"github.com/stretchr/testify/assert"
)

// DECODED {:type=>"data", :logicalTimestampMs=>1504064040000, :logicalTimestamp=>2017-08-30 03:34:00 +0000, :data=>{3079061720=>691.1, 3553579776=>5828.0, 2479549961=>9939.4, 2038453579=>94.8, 3928812177=>2952.2, 2885058095=>686.0, 3689271047=>695.4, 683255203=>756.3, 202128297=>5800.4, 1462695611=>2796.9, 3391947226=>44.8, 1321572762=>1302.3, 1136315563=>8700.8, 122567741=>16128.1, 800290351=>762.5, 1533912710=>2439.7}, :channel=>"channel-1"}
const binaryMsgBase64 = "AQUBAGNoYW5uZWwtMQAAAAAAAAAfiwgAAAAAAAAAY2BgjDM0fuPAwMAgwAQkGLa3Hbjh0DrjDBCcBQtcPh/M4LDtCIjJABaYfOodp8PhzM3GQAAWqGwI9XYIB/MhAq90Xk10WC6QBgJggdXf6vQdWgsQZtx+ep7doXU3QovGzqWLHdoXIbTw8FivdNi2AiEQrvdtt8PSmwiHndJlvOXgloZQ4Xd05iyHKZEIQ513nVntcOAfQgW7n42tw/kGHrgZ+pvL9B3aryAcFp27rM1hMT9YCwAua3WrHAEAAA=="

func TestDecodeBinaryMessage(t *testing.T) {
	rawMsg, err := base64.StdEncoding.DecodeString(binaryMsgBase64)
	if err != nil {
		panic("Could not decode test message")
	}

	msg, err := parseBinaryMessage(rawMsg)
	if err != nil {
		t.Fatalf("could not parse message: %v", err)
	}

	dm, ok := msg.(*DataMessage)
	assert.True(t, ok, "message was not data message")
	assert.NotNil(t, dm)

	assert.Equal(t, dm.Type(), DataType)
	assert.Equal(t, dm.Channel(), "channel-1")
	assert.Equal(t, dm.TimestampMillis, uint64(1504064040000))
	assert.Equal(t, dm.Timestamp().Unix(), time.Date(2017, 8, 29, 23, 34, 00, 0, time.FixedZone("EDT", -4*60*60)).Unix())
	assert.Len(t, dm.Payloads, 16)
	assert.Equal(t, dm.Payloads[0].Value(), 691.1)
	assert.Equal(t, dm.Payloads[0].TSID, idtool.ID(3079061720))
}
