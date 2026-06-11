package rntbd

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"net/http/httptest"

	"github.com/pikami/cosmium/api"
	"github.com/pikami/cosmium/internal/logger"
	tlsprovider "github.com/pikami/cosmium/internal/tls_provider"
)

type RntbdServer struct {
	port      int
	listener  net.Listener
	apiServer *api.ApiServer
}

func NewRntbdServer(port int, apiServer *api.ApiServer) *RntbdServer {
	return &RntbdServer{port: port, apiServer: apiServer}
}

func (s *RntbdServer) Start() error {
	tlsConfig := tlsprovider.GetDefaultTlsConfig()
	listener, err := tls.Listen("tcp", fmt.Sprintf(":%d", s.port), tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %w", s.port, err)
	}
	s.listener = listener

	go func() {
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				logger.ErrorLn("Failed to accept connection:", err)
				continue
			}

			go s.handleConnection(conn)
		}
	}()

	return nil
}

func (s *RntbdServer) Stop() error {
	return s.listener.Close()
}

func (s *RntbdServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		_, err := reader.Peek(4)
		if err != nil {
			return
		}

		frame, err := ReadFrame(reader)
		if err != nil {
			logger.ErrorLn("Failed to read frame:", err)
			continue
		}

		if frame.ResourceType == RntbdResourceTypeConnection {
			responseFrame := buildContextFrame(frame)
			_, err := writer.Write(responseFrame)
			writer.Flush()
			if err != nil {
				logger.ErrorLn("Failed to write response frame:", err)
				continue
			}
			continue
		} else if frame.ResourceType == RntbdResourceTypeDatabase ||
			frame.ResourceType == RntbdResourceTypeCollection ||
			frame.ResourceType == RntbdResourceTypeDocument {
			responseFrameBytes := s.passToApiServer(frame)
			_, err := writer.Write(responseFrameBytes)
			writer.Flush()
			if err != nil {
				logger.ErrorLn("Failed to write response frame:", err)
				continue
			}
			continue
		} else {
			logger.Errorf("Received Unhandled RNTBD request from: %s with resource type: %s\n", conn.RemoteAddr(), frame.ResourceType.String())
		}
	}
}

func (s *RntbdServer) passToApiServer(frame *RntbdFrame) []byte {
	req := frame.ToHttpRequest()
	responseWriter := httptest.NewRecorder()
	s.apiServer.GetRouter().ServeHTTP(responseWriter, req)

	responseFrameBuilder := ToRntbdResponseFrame(responseWriter)
	responseFrameBuilder.SetActivityId(frame.ActivityId)
	if transportRequestId, ok := frame.RequestHeaders[RntbdRequestHeaderTransportRequestID]; ok {
		responseFrameBuilder.AddHeader(uint16(RntbdResponseHeaderTransportRequestID), RntbdTokenTypeULong, transportRequestId)
	}

	responseFrameBuilder.AddHeader(uint16(RntbdResponseHeaderItemLSN), RntbdTokenTypeLongLong, int64(420))
	responseFrameBuilder.AddHeader(uint16(RntbdResponseHeaderLocalLSN), RntbdTokenTypeLongLong, int64(420))
	responseFrameBuilder.AddHeader(uint16(RntbdResponseHeaderGlobalCommittedLSN), RntbdTokenTypeLongLong, int64(420))
	responseFrameBuilder.AddHeader(uint16(RntbdResponseHeaderItemLocalLSN), RntbdTokenTypeLongLong, int64(420))
	responseFrameBuilder.AddHeader(uint16(RntbdResponseHeaderLSN), RntbdTokenTypeLongLong, int64(420))
	responseFrameBuilder.AddHeader(uint16(RntbdResponseHeaderQuorumAckedLSN), RntbdTokenTypeLongLong, int64(420))
	responseFrameBuilder.AddHeader(uint16(RntbdResponseHeaderQuorumAckedLocalLSN), RntbdTokenTypeLongLong, int64(420))
	responseFrameBuilder.AddHeader(uint16(RntbdResponseHeaderCurrentReplicaSetSize), RntbdTokenTypeULong, uint32(1))

	responseFrame := responseFrameBuilder.Build()
	responseFrameBytes := responseFrame.ToBytes()

	return responseFrameBytes
}
