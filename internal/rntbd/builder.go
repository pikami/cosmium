package rntbd

import (
	"bytes"
	"encoding/binary"
)

type RntbdResponseFrame struct {
	StatusCode      uint16
	ResourceType    RntbdResourceType
	ActivityId      []byte
	ResponseHeaders []RntbdResponseHeader
	Payload         []byte
}

type RntbdResponseHeader struct {
	HeaderId   uint16
	TokenType  RntbdTokenType
	TokenValue any
}

type RntbdResponseFrameBuilder struct {
	frame RntbdResponseFrame
}

func (b *RntbdResponseFrameBuilder) AddHeader(headerId uint16, tokenType RntbdTokenType, tokenValue any) {
	b.frame.ResponseHeaders = append(b.frame.ResponseHeaders, RntbdResponseHeader{
		HeaderId:   headerId,
		TokenType:  tokenType,
		TokenValue: tokenValue,
	})
}

func (b *RntbdResponseFrameBuilder) AddPayload(payload []byte) {
	b.frame.Payload = payload
}

func (b *RntbdResponseFrameBuilder) SetStatusCode(statusCode uint16) {
	b.frame.StatusCode = statusCode
}

func (b *RntbdResponseFrameBuilder) SetResourceType(resourceType RntbdResourceType) {
	b.frame.ResourceType = resourceType
}

func (b *RntbdResponseFrameBuilder) SetActivityId(activityId []byte) {
	b.frame.ActivityId = activityId
}

func (b *RntbdResponseFrameBuilder) Build() *RntbdResponseFrame {
	return &b.frame
}

func (f *RntbdResponseFrame) ToBytes() []byte {
	var buffer bytes.Buffer

	binary.Write(&buffer, binary.LittleEndian, f.StatusCode)
	binary.Write(&buffer, binary.LittleEndian, uint16(f.ResourceType))
	binary.Write(&buffer, binary.LittleEndian, f.ActivityId)

	for _, header := range f.ResponseHeaders {
		binary.Write(&buffer, binary.LittleEndian, header.HeaderId)
		binary.Write(&buffer, binary.LittleEndian, uint8(header.TokenType))

		switch header.TokenType {
		case RntbdTokenTypeByte:
			buffer.Write(header.TokenValue.([]byte))
		case RntbdTokenTypeUShort:
			binary.Write(&buffer, binary.LittleEndian, header.TokenValue.(uint16))
		case RntbdTokenTypeULong:
			binary.Write(&buffer, binary.LittleEndian, header.TokenValue.(uint32))
		case RntbdTokenTypeLong:
			binary.Write(&buffer, binary.LittleEndian, header.TokenValue.(int32))
		case RntbdTokenTypeULongLong:
			binary.Write(&buffer, binary.LittleEndian, header.TokenValue.(uint64))
		case RntbdTokenTypeLongLong:
			binary.Write(&buffer, binary.LittleEndian, header.TokenValue.(int64))
		case RntbdTokenTypeGuid:
			buffer.Write(header.TokenValue.([]byte))

		case RntbdTokenTypeSmallString:
			binary.Write(&buffer, binary.LittleEndian, uint8(len(header.TokenValue.(string))))
			buffer.WriteString(header.TokenValue.(string))
		case RntbdTokenTypeString:
			binary.Write(&buffer, binary.LittleEndian, uint16(len(header.TokenValue.(string))))
			buffer.WriteString(header.TokenValue.(string))
		case RntbdTokenTypeULongString:
			binary.Write(&buffer, binary.LittleEndian, uint32(len(header.TokenValue.(string))))
			buffer.WriteString(header.TokenValue.(string))
		case RntbdTokenTypeSmallBytes:
			binary.Write(&buffer, binary.LittleEndian, uint8(len(header.TokenValue.([]byte))))
			buffer.Write(header.TokenValue.([]byte))
		case RntbdTokenTypeBytes:
			binary.Write(&buffer, binary.LittleEndian, uint16(len(header.TokenValue.([]byte))))
			buffer.Write(header.TokenValue.([]byte))
		case RntbdTokenTypeULongBytes:
			binary.Write(&buffer, binary.LittleEndian, uint32(len(header.TokenValue.([]byte))))
			buffer.Write(header.TokenValue.([]byte))
		case RntbdTokenTypeFloat:
			binary.Write(&buffer, binary.LittleEndian, header.TokenValue.(float32))
		case RntbdTokenTypeDouble:
			binary.Write(&buffer, binary.LittleEndian, header.TokenValue.(float64))
		case RntbdTokenTypeInvalid:
			panic("invalid token type")
		default:
			panic("invalid token type")
		}
	}

	payloadSize := uint32(0)
	if len(f.Payload) > 0 {
		payloadSize = uint32(len(f.Payload)) + 4
	}

	frameSize := uint32(buffer.Len()) + 4
	result := make([]byte, frameSize+payloadSize)
	binary.LittleEndian.PutUint32(result, frameSize)
	copy(result[4:], buffer.Bytes())

	if len(f.Payload) > 0 {
		binary.LittleEndian.PutUint32(result[frameSize:], payloadSize-4)
		copy(result[frameSize+4:], f.Payload)
	}

	return result
}

func buildContextFrame(requestFrame *RntbdFrame) []byte {
	builder := RntbdResponseFrameBuilder{}
	builder.SetStatusCode(200)
	builder.SetResourceType(RntbdResourceTypeConnection)
	builder.SetActivityId(requestFrame.ActivityId)
	builder.AddHeader(uint16(RntbdContextHeaderServerAgent), RntbdTokenTypeSmallString, "DocumentDB Server")
	builder.AddHeader(uint16(RntbdContextHeaderServerVersion), RntbdTokenTypeSmallString, " version=2.14.0.0")
	builder.AddHeader(uint16(RntbdContextHeaderIdleTimeoutInSeconds), RntbdTokenTypeULong, uint32(120))
	builder.AddHeader(uint16(RntbdContextHeaderUnauthenticatedTimeoutInSeconds), RntbdTokenTypeULong, uint32(25))
	return builder.Build().ToBytes()
}
