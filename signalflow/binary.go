package signalflow

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"io/ioutil"
)

type binaryMessageHeader struct {
	TimestampMillis uint64
	ElementCount    uint32
}

const (
	ValTypeLong   uint8 = 1
	ValTypeDouble uint8 = 2
	ValTypeInt    uint8 = 3
)

// The first 20 bytes of every binary websocket message from the backend.
// https://developers.signalfx.com/signalflow_analytics/rest_api_messages/stream_messages_specification.html#_binary_encoding_of_websocket_messages
type msgHeader struct {
	Version     uint8
	MessageType uint8
	Flags       uint8
	Reserved    uint8
	Channel     [16]byte
}

const (
	compressed  uint8 = 1 << iota
	jsonEncoded       = 1 << iota
)

func parseBinaryHeader(msg []byte) (string, bool /* isCompressed */, bool /* isJSON */, []byte /* rest of message */, error) {
	if len(msg) <= 20 {
		return "", false, false, nil, fmt.Errorf("expected SignalFlow message of at least 21 bytes, got %d bytes", len(msg))
	}

	r := bytes.NewReader(msg[:20])
	var header msgHeader
	err := binary.Read(r, binary.BigEndian, &header)
	if err != nil {
		return "", false, false, nil, err
	}

	isCompressed := header.Flags&compressed != 0
	isJSON := header.Flags&jsonEncoded != 0

	return string(header.Channel[:bytes.IndexByte(header.Channel[:], 0)]), isCompressed, isJSON, msg[20:], err
}

func parseBinaryMessage(msg []byte) (Message, error) {
	channel, isCompressed, isJSON, rest, err := parseBinaryHeader(msg)
	if err != nil {
		return nil, err
	}

	if isCompressed {
		var err error
		reader, err := gzip.NewReader(bytes.NewReader(rest))
		if err != nil {
			return nil, err
		}
		rest, err = ioutil.ReadAll(reader)
		if err != nil {
			return nil, err
		}
	}

	if isJSON {
		panic("cannot handle json binary message")
	}

	r := bytes.NewReader(rest[:12])
	var header binaryMessageHeader
	err = binary.Read(r, binary.BigEndian, &header)
	if err != nil {
		return nil, err
	}

	var payloads []BinaryPayload
	for i := 0; i < int(header.ElementCount); i++ {
		r := bytes.NewReader(rest[12+17*i : 12+17*(i+1)])
		var payload BinaryPayload
		if err := binary.Read(r, binary.BigEndian, &payload); err != nil {
			return nil, err
		}
		payloads = append(payloads, payload)
	}

	return &DataMessage{
		BaseMessage: BaseMessage{
			Typ:  DataType,
			Chan: channel,
		},
		TimestampMillis: header.TimestampMillis,
		Payloads:        payloads,
	}, nil
}
