package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/golang/protobuf/jsonpb"
	pb "github.com/kazegusuri/grpcurl/internal/testdata"
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
	expected := `{"value":"xxx","error_code":0}`
	assert.Equal(t, expected, resp.RequestMessage, "request message")
	assert.Equal(t, expected, resp.ResponseMessage, "response message")
}

func TestCallEchoEmpty(t *testing.T) {
	buf, err := testCall("grpcurl.test.Echo.Echo", `{}`)
	require.NoError(t, err)
	resp := parseTestResponse(buf.String())
	expected := `{"value":"","error_code":0}`
	assert.Equal(t, expected, resp.RequestMessage, "request message")
	assert.Equal(t, expected, resp.ResponseMessage, "response message")
}

func TestCallEchoUnknownField(t *testing.T) {
	buf, err := testCall("grpcurl.test.Echo.Echo", `{"xxx": "vvvv"}`)
	require.NoError(t, err)
	resp := parseTestResponse(buf.String())
	expected := `{"value":"","error_code":0}`
	assert.Equal(t, expected, resp.RequestMessage, "request message")
	assert.Equal(t, expected, resp.ResponseMessage, "response message")
}

func TestEverythingSimple(t *testing.T) {
	buf, err := testCall(
		"grpcurl.test.Everything.Simple",
		`{"string_value": "aaa", "bool_value": true}`)
	require.NoError(t, err)
	resp := parseTestResponse(buf.String())
	expected := `{"string_value":"aaa","bool_value":true}`
	assert.Equal(t, expected, resp.RequestMessage, "request message")
	assert.Equal(t, expected, resp.ResponseMessage, "response message")
}

func TestEverythingSimpleEmpty(t *testing.T) {
	buf, err := testCall(
		"grpcurl.test.Everything.Simple",
		`{}`)
	require.NoError(t, err)
	resp := parseTestResponse(buf.String())
	expected := `{"string_value":"","bool_value":false}`
	assert.Equal(t, expected, resp.RequestMessage, "request message")
	assert.Equal(t, expected, resp.ResponseMessage, "response message")
}

func TestEverythingNumber(t *testing.T) {
	m := &jsonpb.Marshaler{OrigName: true}
	msg, err := m.MarshalToString(&pb.NumberMessage{
		FloatValue:    1.1,
		DoubleValue:   2.2,
		Int32Value:    3,
		Int64Value:    4,
		Uint32Value:   5,
		Uint64Value:   6,
		Sint32Value:   7,
		Sint64Value:   8,
		Fixed32Value:  9,
		Fixed64Value:  10,
		Sfixed32Value: 11,
		Sfixed64Value: 12,
	})
	require.NoError(t, err)
	t.Logf("request messge: %q", msg)

	buf, err := testCall("grpcurl.test.Everything.Number", msg)
	require.NoError(t, err)
	resp := parseTestResponse(buf.String())
	expected := `{"float_value":1.1,"double_value":2.2,"int32_value":3,"int64_value":4,"uint32_value":5,"uint64_value":6,"sint32_value":7,"sint64_value":8,"fixed32_value":9,"fixed64_value":10,"sfixed32_value":11,"sfixed64_value":12}`
	assert.Equal(t, expected, resp.RequestMessage, "request message")
	assert.Equal(t, expected, resp.ResponseMessage, "response message")
}

func TestEverythingNumberEmpty(t *testing.T) {
	buf, err := testCall("grpcurl.test.Everything.Number", `{}`)
	require.NoError(t, err)
	resp := parseTestResponse(buf.String())
	expected := `{"float_value":0,"double_value":0,"int32_value":0,"int64_value":0,"uint32_value":0,"uint64_value":0,"sint32_value":0,"sint64_value":0,"fixed32_value":0,"fixed64_value":0,"sfixed32_value":0,"sfixed64_value":0}`
	assert.Equal(t, expected, resp.RequestMessage, "request message")
	assert.Equal(t, expected, resp.ResponseMessage, "response message")
}

func TestEverythingEnum(t *testing.T) {
	m := &jsonpb.Marshaler{OrigName: true}
	msg, err := m.MarshalToString(&pb.EnumMessage{
		NumericEnumValue: pb.NumericEnum_ONE,
		RepeatedNumericEnumValues: []pb.NumericEnum{
			pb.NumericEnum_ONE,
			pb.NumericEnum_TWO,
		},
		AliasedEnumValue: pb.AliasedEnum_RUNNING,
		NestedEnumValue:  pb.EnumMessage_PENDING,
		RepeatedNestedEnumValues: []pb.EnumMessage_Nested{
			pb.EnumMessage_PENDING,
			pb.EnumMessage_COMPLETED,
		},
	})
	require.NoError(t, err)
	t.Logf("request messge: %q", msg)

	buf, err := testCall("grpcurl.test.Everything.Enum", msg)
	require.NoError(t, err)
	resp := parseTestResponse(buf.String())
	expected := `{"numeric_enum_value":"ONE","repeated_numeric_enum_values":["ONE","TWO"],"aliased_enum_value":"STARTED","nested_enum_value":"PENDING","repeated_nested_enum_values":["PENDING","COMPLETED"]}`
	assert.Equal(t, expected, resp.RequestMessage, "request message")
	assert.Equal(t, expected, resp.ResponseMessage, "response message")
}

func TestEverythingEnumAsInt(t *testing.T) {
	m := &jsonpb.Marshaler{EnumsAsInts: true, OrigName: true}
	msg, err := m.MarshalToString(&pb.EnumMessage{
		NumericEnumValue: pb.NumericEnum_ONE,
		RepeatedNumericEnumValues: []pb.NumericEnum{
			pb.NumericEnum_ONE,
			pb.NumericEnum_TWO,
		},
		AliasedEnumValue: pb.AliasedEnum_RUNNING,
		NestedEnumValue:  pb.EnumMessage_PENDING,
		RepeatedNestedEnumValues: []pb.EnumMessage_Nested{
			pb.EnumMessage_PENDING,
			pb.EnumMessage_COMPLETED,
		},
	})
	require.NoError(t, err)
	t.Logf("request messge: %q", msg)

	buf, err := testCall("grpcurl.test.Everything.Enum", msg)
	require.NoError(t, err)
	resp := parseTestResponse(buf.String())
	expected := `{"numeric_enum_value":"ONE","repeated_numeric_enum_values":["ONE","TWO"],"aliased_enum_value":"STARTED","nested_enum_value":"PENDING","repeated_nested_enum_values":["PENDING","COMPLETED"]}`
	assert.Equal(t, expected, resp.RequestMessage, "request message")
	assert.Equal(t, expected, resp.ResponseMessage, "response message")
}

func TestEverythingEnumEmpty(t *testing.T) {
	buf, err := testCall("grpcurl.test.Everything.Enum", `{}`)
	require.NoError(t, err)
	resp := parseTestResponse(buf.String())
	expected := `{"numeric_enum_value":"ZERO","repeated_numeric_enum_values":[],"aliased_enum_value":"UNKNOWN","nested_enum_value":"UNKNOWN","repeated_nested_enum_values":[]}`
	assert.Equal(t, expected, resp.RequestMessage, "request message")
	assert.Equal(t, expected, resp.ResponseMessage, "response message")
}

func TestEverythingOneof1(t *testing.T) {
	buf, err := testCall("grpcurl.test.Everything.Oneof", `{"int32_value": 100}`)
	require.NoError(t, err)
	resp := parseTestResponse(buf.String())
	expected := `{"int32_value":100,"string_value":"","repeated_oneof_values":[]}`
	assert.Equal(t, expected, resp.RequestMessage, "request message")
	assert.Equal(t, expected, resp.ResponseMessage, "response message")
}

func TestEverythingOneof2(t *testing.T) {
	buf, err := testCall("grpcurl.test.Everything.Oneof", `{"string_value": "xxx"}`)
	require.NoError(t, err)
	resp := parseTestResponse(buf.String())
	expected := `{"int32_value":0,"string_value":"xxx","repeated_oneof_values":[]}`
	assert.Equal(t, expected, resp.RequestMessage, "request message")
	assert.Equal(t, expected, resp.ResponseMessage, "response message")
}

func TestEverythingOneof3(t *testing.T) {
	m := &jsonpb.Marshaler{EnumsAsInts: true, OrigName: true}
	msg, err := m.MarshalToString(&pb.OneofMessage{
		OneofValue: &pb.OneofMessage_Int32Value{Int32Value: 2000},
	})
	require.NoError(t, err)
	t.Logf("request messge: %q", msg)

	buf, err := testCall("grpcurl.test.Everything.Oneof", msg)
	require.NoError(t, err)
	resp := parseTestResponse(buf.String())
	expected := `{"int32_value":2000,"string_value":"","repeated_oneof_values":[]}`
	assert.Equal(t, expected, resp.RequestMessage, "request message")
	assert.Equal(t, expected, resp.ResponseMessage, "response message")
}

func TestEverythingOneofEmpty(t *testing.T) {
	buf, err := testCall("grpcurl.test.Everything.Oneof", `{}`)
	require.NoError(t, err)
	resp := parseTestResponse(buf.String())
	expected := `{"int32_value":0,"string_value":"","repeated_oneof_values":[]}`
	assert.Equal(t, expected, resp.RequestMessage, "request message")
	assert.Equal(t, expected, resp.ResponseMessage, "response message")
}

func TestEverythingMap(t *testing.T) {
	m := &jsonpb.Marshaler{OrigName: true}
	msg, err := m.MarshalToString(&pb.MapMessage{
		MappedValue: map[string]string{
			"foo": "foo1",
			"bar": "bar1",
		},
		MappedEnumValue: map[string]pb.NumericEnum{
			"one": pb.NumericEnum_ONE,
			"two": pb.NumericEnum_TWO,
		},
		MappedNestedValue: map[string]*pb.NestedMessage{
			"foo": &pb.NestedMessage{
				NestedValue: &pb.NestedMessage_Nested{
					Int32Value:  100,
					StringValue: "xxx",
				},
				RepeatedNestedValues: []*pb.NestedMessage_Nested{
					{
						Int32Value:  200,
						StringValue: "yyy",
					},
					{
						Int32Value:  300,
						StringValue: "zzz",
					},
				},
			},
		},
	})
	require.NoError(t, err)
	t.Logf("request messge: %q", msg)

	buf, err := testCall("grpcurl.test.Everything.Map", msg)
	require.NoError(t, err)
	resp := parseTestResponse(buf.String())
	expected := `{"mapped_value":{"bar":"bar1","foo":"foo1"},"mapped_enum_value":{"one":"ONE","two":"TWO"},"mapped_nested_value":{"foo":{"nested_value":{"int32_value":100,"string_value":"xxx"},"repeated_nested_values":[{"int32_value":200,"string_value":"yyy"},{"int32_value":300,"string_value":"zzz"}]}}}`
	assert.Equal(t, expected, resp.RequestMessage, "request message")
	assert.Equal(t, expected, resp.ResponseMessage, "response message")
}

// TODO protoreflect support
func TestEverythingMapEmpty(t *testing.T) {
	t.Skip("protoreflect not supported")
	buf, err := testCall("grpcurl.test.Everything.Map", `{}`)
	require.NoError(t, err)
	resp := parseTestResponse(buf.String())
	expected := `{}`
	assert.Equal(t, expected, resp.RequestMessage, "request message")
	assert.Equal(t, expected, resp.ResponseMessage, "response message")
}

func TestEchoV2Echo(t *testing.T) {
	// cross package message
	buf, err := testCall(
		"grpcurl.test.v2.Echo.Echo",
		`{"value": "xxx"}`)
	require.NoError(t, err)
	resp := parseTestResponse(buf.String())
	expected := `{"value":"xxx","error_code":0}`
	assert.Equal(t, expected, resp.RequestMessage, "request message")
	assert.Equal(t, expected, resp.ResponseMessage, "response message")
}
