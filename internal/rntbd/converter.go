package rntbd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"

	"github.com/pikami/cosmium/api/headers"
)

func (f *RntbdFrame) ToHttpRequest() *http.Request {
	req := &http.Request{
		Method: operationTypeToHttpMethod(f.OperationType),
		URL:    &url.URL{Path: frameToPath(f)},
		Body:   io.NopCloser(bytes.NewReader(f.Payload)),
		Header: http.Header{},
	}

	switch f.OperationType {
	case RntbdOperationTypeQuery, RntbdOperationTypeSQLQuery:
		req.Header.Set(headers.Query, "true")
	case RntbdOperationTypeUpsert:
		req.Header.Set(headers.IsUpsert, "true")
	}

	if ifMatch, ok := f.RequestHeaders[RntbdRequestHeaderMatch]; ok {
		if ifMatchString, ok := ifMatch.(string); ok {
			req.Header.Set(headers.IfMatch, ifMatchString)
		}
	}

	if continuationToken, ok := f.RequestHeaders[RntbdRequestHeaderContinuationToken]; ok {
		if continuationTokenString, ok := continuationToken.(string); ok {
			req.Header.Set(headers.ContinuationToken, continuationTokenString)
		}
	}

	if maxItemCount, ok := f.RequestHeaders[RntbdRequestHeaderPageSize]; ok {
		if maxItemCountString, ok := maxItemCount.(uint64); ok {
			req.Header.Set(headers.MaxItemCount, fmt.Sprintf("%d", maxItemCountString))
		}
	}

	return req
}

func ToRntbdResponseFrame(responseWriter *httptest.ResponseRecorder) *RntbdResponseFrameBuilder {
	builder := &RntbdResponseFrameBuilder{}
	builder.SetStatusCode(uint16(responseWriter.Code))

	if responseWriter.Header().Get(headers.ETag) != "" {
		builder.AddHeader(uint16(RntbdResponseHeaderETag), RntbdTokenTypeString, responseWriter.Header().Get(headers.ETag))
	}

	if responseWriter.Header().Get(headers.ContinuationToken) != "" {
		builder.AddHeader(uint16(RntbdResponseHeaderContinuationToken), RntbdTokenTypeString, responseWriter.Header().Get(headers.ContinuationToken))
	}

	if responseWriter.Header().Get(headers.ItemCount) != "" {
		itemCount, err := strconv.ParseUint(responseWriter.Header().Get(headers.ItemCount), 10, 32)
		if err != nil {
			panic(err)
		}

		builder.AddHeader(uint16(RntbdResponseHeaderItemCount), RntbdTokenTypeULong, uint32(itemCount))
	}

	if responseWriter.Body.Len() > 0 {
		builder.AddHeader(uint16(RntbdResponseHeaderPayloadPresent), RntbdTokenTypeByte, []byte{1})
		builder.AddPayload(responseWriter.Body.Bytes())
	} else {
		builder.AddHeader(uint16(RntbdResponseHeaderPayloadPresent), RntbdTokenTypeByte, []byte{0})
	}

	return builder
}

func operationTypeToHttpMethod(operationType RntbdOperationType) string {
	switch operationType {
	case RntbdOperationTypeRead,
		RntbdOperationTypeReadFeed:
		return http.MethodGet
	case RntbdOperationTypeCreate,
		RntbdOperationTypeUpsert,
		RntbdOperationTypeQuery,
		RntbdOperationTypeSQLQuery:
		return http.MethodPost
	case RntbdOperationTypeUpdate,
		RntbdOperationTypeReplace:
		return http.MethodPut
	case RntbdOperationTypeDelete:
		return http.MethodDelete
	}

	panic(fmt.Sprintf("Unknown operation type: %d", operationType))
}

func frameToPath(frame *RntbdFrame) string {
	databaseName, databaseOk := frame.RequestHeaders[RntbdRequestHeaderDatabaseName]
	collectionName, collectionOk := frame.RequestHeaders[RntbdRequestHeaderCollectionName]
	documentName, documentOk := frame.RequestHeaders[RntbdRequestHeaderDocumentName]

	urlPath := ""

	if databaseOk {
		urlPath += fmt.Sprintf("/dbs/%s", databaseName)
	} else if frame.ResourceType == RntbdResourceTypeDatabase {
		urlPath += "/dbs"
	}

	if collectionOk {
		urlPath += fmt.Sprintf("/colls/%s", collectionName)
	} else if frame.ResourceType == RntbdResourceTypeCollection {
		urlPath += "/colls"
	}

	if documentOk {
		urlPath += fmt.Sprintf("/docs/%s", documentName)
	} else if frame.ResourceType == RntbdResourceTypeDocument {
		urlPath += "/docs"
	}

	return urlPath
}
