package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ExampleCall() {
	cmd := NewRootCommand(strings.NewReader(`{"value": "hello"}`), os.Stdout)
	cmd.Command().SetArgs([]string{"-k", "call", addr, "grpcurl.test.Echo.Echo"})
	cmd.Command().Execute()
	// Output:
	// {"value":"hello","error_code":0}
}

func testCall(method, msg string) (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}
	buf.Grow(1024)
	cmd := NewRootCommand(strings.NewReader(msg), buf)
	cmd.Command().SetArgs([]string{"-k", "call", "-v", addr, method})
	return buf, cmd.Command().Execute()
}

type testResponse struct {
	RequestMessage  string
	ResponseMessage string
	RequestHeader   map[string][]string
	ResponseHeader  map[string][]string
	ResponseTrailer map[string][]string
}

func parseTestResponse(s string) *testResponse {
	lines := strings.Split(s, "\n")

	msgs := map[string][]string{}
	marker := ""
	for i := range lines {
		switch lines[i] {
		case requestMessageMarker:
			marker = requestMessageMarker
		case responseMessageMarker:
			marker = responseMessageMarker
		case responseHeaderMarker:
			marker = responseHeaderMarker
		case responseTrailerMarker:
			marker = responseTrailerMarker
		default:
			msgs[marker] = append(msgs[marker], lines[i])
		}
	}

	return &testResponse{
		RequestMessage:  strings.Join(msgs[requestMessageMarker], "\n"),
		ResponseMessage: strings.Join(msgs[responseMessageMarker], "\n"),
	}
}

func TestCallEcho(t *testing.T) {
	buf, err := testCall("grpcurl.test.Echo.Echo", `{"value": "xxx"}`)
	require.NoError(t, err)
	resp := parseTestResponse(buf.String())
	assert.Equal(t, resp.RequestMessage, `{"value":"xxx","error_code":0}`)
	assert.Equal(t, resp.RequestMessage, `{"value":"xxx","error_code":0}`)
}

func TestCallEchoEmpty(t *testing.T) {
	buf, err := testCall("grpcurl.test.Echo.Echo", `{}`)
	require.NoError(t, err)
	resp := parseTestResponse(buf.String())
	assert.Equal(t, resp.RequestMessage, `{"value":"","error_code":0}`)
	assert.Equal(t, resp.RequestMessage, `{"value":"","error_code":0}`)
}

func TestCallEchoUnknownField(t *testing.T) {
	buf, err := testCall("grpcurl.test.Echo.Echo", `{"xxx": "vvvv"}`)
	require.NoError(t, err)
	resp := parseTestResponse(buf.String())
	assert.Equal(t, resp.RequestMessage, `{"value":"","error_code":0}`)
	assert.Equal(t, resp.RequestMessage, `{"value":"","error_code":0}`)
}
