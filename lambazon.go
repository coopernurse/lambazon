package lambazon

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	apex "github.com/apex/go-apex"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

/////////////////////////

type ClosingBuffer struct {
	*bytes.Buffer
}

func (cb *ClosingBuffer) Close() error {
	// no-op
	return nil
}

func NewResponseWriter() *ResponseWriter {
	return &ResponseWriter{
		status: 200,
		buf:    &bytes.Buffer{},
		header: http.Header{},
	}
}

type ResponseWriter struct {
	status int
	buf    *bytes.Buffer
	header http.Header
}

func (w *ResponseWriter) Header() http.Header {
	return w.header
}

func (w *ResponseWriter) Write(data []byte) (int, error) {
	return w.buf.Write(data)
}

func (w *ResponseWriter) WriteHeader(status int) {
	w.status = status
}

func (w *ResponseWriter) ToReply() Reply {
	var body string
	var bodyEncoding string
	if w.shouldBase64Encode() {
		body = base64.StdEncoding.EncodeToString(w.buf.Bytes())
		bodyEncoding = "base64"
	} else {
		body = w.buf.String()
	}

	return Reply{
		Type: "HTTPJSON-REP",
		Meta: &ReplyMeta{
			Status:  w.status,
			Headers: w.header,
		},
		Body:         body,
		BodyEncoding: bodyEncoding,
	}
}

func (w *ResponseWriter) shouldBase64Encode() bool {
	if w.buf.Len() == 0 {
		return false
	}

	ct := w.header.Get("content-type")
	if ct != "" && ct != "application/json" && !strings.HasPrefix(ct, "text/") {
		return true
	}

	return false
}

/////////////////////////

// Converts the inbound Lambda JSON event to a net/http Request
func toRequest(event json.RawMessage) (*http.Request, error) {
	var reqEnvelope Request
	if err := json.Unmarshal(event, &reqEnvelope); err != nil {
		return nil, err
	}
	if reqEnvelope.Meta == nil || reqEnvelope.Meta.Method == "" {
		return nil, fmt.Errorf("Request meta not provided")
	}

	req := &http.Request{}
	req.Method = reqEnvelope.Meta.Method

	path := reqEnvelope.Meta.Path
	if reqEnvelope.Meta.Query != "" {
		path += "?" + reqEnvelope.Meta.Query
	}
	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	req.URL = u

	req.Body = &ClosingBuffer{bytes.NewBufferString(reqEnvelope.Body)}
	req.Proto = reqEnvelope.Meta.Proto
	req.Host = reqEnvelope.Meta.Host

	req.Header = make(http.Header)
	for k, xs := range reqEnvelope.Meta.Headers {
		for _, s := range xs {
			req.Header.Add(k, s)
		}
	}

	if len(reqEnvelope.Body) > 0 {
		req.Header.Set("Content-Length", strconv.Itoa(len(reqEnvelope.Body)))
	}

	return req, nil
}

// Request represents a single HTTP request.  It will be serialized as JSON
// and sent to the AWS Lambda function as the function payload.
type Request struct {
	// Set to the constant "HTTPJSON-REQ"
	Type string `json:"type"`
	// Metadata about the HTTP request
	Meta *RequestMeta `json:"meta"`
	// HTTP request body (may be empty)
	Body string `json:"body"`
}

// RequestMeta represents HTTP metadata present on the request
type RequestMeta struct {
	// HTTP method used by client (e.g. GET or POST)
	Method string `json:"method"`

	// Path portion of URL without the query string
	Path string `json:"path"`

	// Query string (without '?')
	Query string `json:"query"`

	// Host field from net/http Request, which may be of the form host:port
	Host string `json:"host"`

	// Proto field from net/http Request, for example "HTTP/1.1"
	Proto string `json:"proto"`

	// HTTP request headers
	Headers map[string][]string `json:"headers"`
}

type Reply struct {
	// Must be set to the constant "HTTPJSON-REP"
	Type string `json:"type"`
	// Reply metadata. If omitted, a default 200 status with empty headers will be used.
	Meta *ReplyMeta `json:"meta"`
	// Response body
	Body string `json:"body"`
	// Encoding of Body - Valid values: "", "base64"
	BodyEncoding string `json:"bodyEncoding"`
}

// ReplyMeta encapsulates HTTP response metadata that the lambda function wishes
// Caddy to set on the HTTP response.
//
// *NOTE* that header values must be encoded as string arrays
type ReplyMeta struct {
	// HTTP status code (e.g. 200 or 404)
	Status int `json:"status"`
	// HTTP response headers
	Headers map[string][]string `json:"headers"`
}

/////////////////////////////

func Run(h http.Handler) {
	apex.HandleFunc(func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
		req, err := toRequest(event)
		if err != nil {
			return nil, err
		}
		rw := NewResponseWriter()
		h.ServeHTTP(rw, req)
		return rw.ToReply(), nil
	})
}
