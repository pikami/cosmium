package rntbd

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/pikami/cosmium/internal/logger"
)

type RntbdFrame struct {
	ResourceType    RntbdResourceType
	OperationType   RntbdOperationType
	ActivityId      []byte
	RequestHeaders  map[RntbdRequestHeader]any
	ResponseHeaders map[RntbdResponseHeaderType]any
	ContextHeaders  map[RntbdContextHeader]any
	Payload         []byte
}

func ReadFrame(reader *bufio.Reader) (*RntbdFrame, error) {
	sizeBytes := readBytes(reader, 4)
	size := binary.LittleEndian.Uint32(sizeBytes)

	payload := readBytes(reader, int(size)-4)

	frame, err := parseFrame_Int(payload, false)
	if err != nil {
		return nil, err
	}

	if payloadPresent, ok := frame.RequestHeaders[RntbdRequestHeaderPayloadPresent]; ok && payloadPresent.([]byte)[0] == 1 {
		payloadSize := binary.LittleEndian.Uint32(readBytes(reader, 4))
		payload := readBytes(reader, int(payloadSize))
		frame.Payload = payload
	}

	if payloadPresent, ok := frame.ResponseHeaders[RntbdResponseHeaderPayloadPresent]; ok && payloadPresent.([]byte)[0] == 1 {
		payloadSize := binary.LittleEndian.Uint32(readBytes(reader, 4))
		payload := readBytes(reader, int(payloadSize))
		frame.Payload = payload
	}

	return frame, nil
}

func ParseFrame(data []byte, isResponse bool) (*RntbdFrame, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("data too short")
	}

	reader := bufio.NewReader(bytes.NewReader(data))
	sizeBytes := readBytes(reader, 4)
	size := binary.LittleEndian.Uint32(sizeBytes)

	payload := readBytes(reader, int(size)-4)
	frame, err := parseFrame_Int(payload, isResponse)
	if err != nil {
		return nil, err
	}

	if payloadPresent, ok := frame.RequestHeaders[RntbdRequestHeaderPayloadPresent]; ok && payloadPresent.([]byte)[0] == 1 {
		payloadSize := binary.LittleEndian.Uint32(readBytes(reader, 4))
		payload := readBytes(reader, int(payloadSize))
		frame.Payload = payload
	}

	if payloadPresent, ok := frame.ResponseHeaders[RntbdResponseHeaderPayloadPresent]; ok && payloadPresent.([]byte)[0] == 1 {
		payloadSize := binary.LittleEndian.Uint32(readBytes(reader, 4))
		payload := readBytes(reader, int(payloadSize))
		frame.Payload = payload
	}

	leftOverBytes, err := io.ReadAll(reader)
	if err != nil {
		logger.ErrorLn("Error reading leftOverBytes:", err)
	}

	if len(leftOverBytes) > 0 {
		logger.ErrorLn("Left over bytes:", hex.EncodeToString(leftOverBytes))
	}

	return frame, nil
}

func parseFrame_Int(data []byte, isResponse bool) (*RntbdFrame, error) {
	payloadReader := bufio.NewReader(bytes.NewReader(data))

	resourceTypeBytes := readBytes(payloadReader, 2)
	resourceType := binary.LittleEndian.Uint16(resourceTypeBytes)

	operationTypeBytes := readBytes(payloadReader, 2)
	operationType := RntbdOperationType(binary.LittleEndian.Uint16(operationTypeBytes))

	activityIdBytes := readBytes(payloadReader, 16)

	requestHeaders := make(map[RntbdRequestHeader]any)
	responseHeaders := make(map[RntbdResponseHeaderType]any)
	contextHeaders := make(map[RntbdContextHeader]any)
	for {
		if _, err := payloadReader.Peek(1); err != nil {
			break
		}

		headerIdBytes := readBytes(payloadReader, 2)
		headerId := binary.LittleEndian.Uint16(headerIdBytes)

		token, err := parseRntbdToken(payloadReader)
		if err != nil {
			return nil, err
		}

		if resourceType == uint16(RntbdResourceTypeConnection) {
			contextHeaders[RntbdContextHeader(headerId)] = token
		} else if isResponse {
			responseHeaders[RntbdResponseHeaderType(headerId)] = token
		} else {
			requestHeaders[RntbdRequestHeader(headerId)] = token
		}
	}

	return &RntbdFrame{
		ResourceType:    RntbdResourceType(resourceType),
		OperationType:   RntbdOperationType(operationType),
		ActivityId:      activityIdBytes,
		RequestHeaders:  requestHeaders,
		ResponseHeaders: responseHeaders,
		ContextHeaders:  contextHeaders,
	}, nil
}

func parseRntbdToken(reader *bufio.Reader) (any, error) {
	tokenTypeBytes := readBytes(reader, 1)
	tokenType := RntbdTokenType(tokenTypeBytes[0])

	switch tokenType {
	case RntbdTokenTypeByte:
		token := readBytes(reader, 1)
		return token, nil
	case RntbdTokenTypeUShort:
		token := binary.LittleEndian.Uint16(readBytes(reader, 2))
		return token, nil
	case RntbdTokenTypeULong:
		token := binary.LittleEndian.Uint32(readBytes(reader, 4))
		return token, nil
	case RntbdTokenTypeLong:
		token := int32(binary.LittleEndian.Uint32(readBytes(reader, 4)))
		return token, nil
	case RntbdTokenTypeULongLong:
		token := binary.LittleEndian.Uint64(readBytes(reader, 8))
		return token, nil
	case RntbdTokenTypeLongLong:
		token := int64(binary.LittleEndian.Uint64(readBytes(reader, 8)))
		return token, nil
	case RntbdTokenTypeGuid:
		token := readBytes(reader, 16)
		return token, nil
	case RntbdTokenTypeSmallString:
		lengthBytes := readBytes(reader, 1)
		length := uint8(lengthBytes[0])
		token := readBytes(reader, int(length))
		return string(token), nil
	case RntbdTokenTypeString:
		length := binary.LittleEndian.Uint16(readBytes(reader, 2))
		token := readBytes(reader, int(length))
		return string(token), nil
	case RntbdTokenTypeULongString:
		length := binary.LittleEndian.Uint32(readBytes(reader, 4))
		token := readBytes(reader, int(length))
		return string(token), nil
	case RntbdTokenTypeSmallBytes:
		lengthBytes := readBytes(reader, 1)
		length := uint8(lengthBytes[0])
		token := readBytes(reader, int(length))
		return token, nil
	case RntbdTokenTypeBytes:
		length := binary.LittleEndian.Uint16(readBytes(reader, 2))
		token := readBytes(reader, int(length))
		return token, nil
	case RntbdTokenTypeULongBytes:
		length := binary.LittleEndian.Uint32(readBytes(reader, 4))
		token := readBytes(reader, int(length))
		return token, nil
	case RntbdTokenTypeFloat:
		// I can't be bothered to implement this, let's just return a byte array
		token := readBytes(reader, 4)
		return token, nil
	case RntbdTokenTypeDouble:
		// I can't be bothered to implement this, let's just return a byte array
		token := readBytes(reader, 8)
		return token, nil
	case RntbdTokenTypeInvalid:
		return nil, fmt.Errorf("invalid token type")
	}

	return nil, fmt.Errorf("invalid token type")
}

func readBytes(reader *bufio.Reader, n int) []byte {
	bytes := make([]byte, n)
	_, err := io.ReadFull(reader, bytes)
	if err != nil {
		logger.ErrorLn("Error reading bytes:", err)
		os.Exit(0)
	}

	return bytes
}
