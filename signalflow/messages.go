package signalflow

import (
	"encoding/binary"
	"math"
)

const (
	AuthenticatedType  = "authenticated"
	ControlMessageType = "control-message"
	MetadataType       = "metadata"
	MessageType        = "message"
	DataType           = "data"
	WebsocketErrorType = "websocket-error"
)

const (
	StreamStartEvent  = "STREAM_START"
	JobStartEvent     = "JOB_START"
	JobProgressEvent  = "JOB_PROGRESS"
	ChannelAbortEvent = "CHANNEL_ABORT"
	EndOfChannelEvent = "END_OF_CHANNEL"
)

type BaseMessage struct {
	Typ  string `json:"type"`
	Chan string `json:"channel,omitempty"`
}

func (bm *BaseMessage) Type() string {
	return bm.Typ
}

func (bm *BaseMessage) Channel() string {
	return bm.Chan
}

var _ Message = &BaseMessage{}

type Message interface {
	Type() string
}

type ChannelMessage interface {
	Message
	Channel() string
}

type AuthenticatedMessage struct {
	OrgID  string `json:"orgId"`
	UserID string `json:"userId"`
}

func (am *AuthenticatedMessage) Type() string {
	return AuthenticatedType
}

type BaseControlMessage struct {
	BaseMessage
	Event           string `json:"event"`
	TimestampMillis uint64 `json:"timestampMs"`
}

func (cm *BaseControlMessage) Type() string {
	return ControlMessageType
}

type JobStartControlMessage struct {
	BaseControlMessage
	Handle string `json:"handle"`
}

type MetadataMessage struct {
	BaseMessage
	Rest map[string]interface{} `json:",inline"`
}

func (mm *MetadataMessage) Type() string {
	return MetadataType
}

type MessageBlock struct {
	TimestampMillis    uint64                 `json:"timestampMs"`
	Code               string                 `json:"messageCode"`
	Level              string                 `json:"messageLevel"`
	NumInputTimeseries int                    `json:"numInputTimeSeries"`
	Contents           map[string]interface{} `json:"contents"`
}

type MessageMessage struct {
	BaseMessage
	LogicalTimestampMillis uint64 `json:"logicalTimestampMs"`
	MessageBlock           `json:"message"`
}

func (mm *MessageMessage) Type() string {
	return MessageType
}

type BinaryPayload struct {
	Type uint8
	Tsid int64
	Val  [8]byte
}

// Value returns the numeric value as an interface{}.
func (bp *BinaryPayload) Value() interface{} {
	switch bp.Type {
	case ValTypeLong:
		return bp.Int64()
	case ValTypeDouble:
		return bp.Float64()
	case ValTypeInt:
		return bp.Int32()
	default:
		return nil
	}
}

func (bp *BinaryPayload) Int64() int64 {
	i, _ := binary.Varint(bp.Val[:])
	return i
}

func (bp *BinaryPayload) Float64() float64 {
	bits := binary.BigEndian.Uint64(bp.Val[:])
	return math.Float64frombits(bits)
}

func (bp *BinaryPayload) Int32() int32 {
	// If the value is an int32, the value from bp.Int64 can always be casted
	// without loss.
	return int32(bp.Int64())
}

// DataMessage is a set of datapoints that share a common timestamp
type DataMessage struct {
	BaseMessage
	TimestampMillis uint64
	Payloads        []BinaryPayload
}

// WebsocketErrorMessage is a client-specific error that indicates something
// went wrong with the underlying connection.
type WebsocketErrorMessage struct {
	Error error
}

func (wem *WebsocketErrorMessage) Type() string {
	return WebsocketErrorType
}
