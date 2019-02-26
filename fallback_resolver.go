package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
)

// DynamicAnyResolver is like the default jsonpb.AnyResolver but, instead of
// returning error, it will fallback to an Empty type if the type is not
// registered. This allows the jsonpb Marshaler to not simply give up in case
// an unknown type is encountered. Because the type is Empty, all fields in
// that the unknown type will implicitly end up as unrecognized, and the JSON
// marshaler will emit them as map key/values.
type DynamicAnyResolver struct {
	jsonpb.AnyResolver // https://godoc.org/github.com/golang/protobuf/jsonpb#AnyResolver
}

// Resolve implements jsonpb.AnyResolver.Resolve
func (DynamicAnyResolver) Resolve(typeURL string) (proto.Message, error) {
	msg, err := defaultResolveAny(typeURL)
	if err == nil {
		return msg, nil
	}
	return &empty.Empty{}, nil
}

// copied from https://github.com/golang/protobuf/blob/c823c79ea1570fb5ff454033735a8e68575d1d0f/jsonpb/jsonpb.go#L92-L103
func defaultResolveAny(typeURL string) (proto.Message, error) {
	// Only the part of typeUrl after the last slash is relevant.
	mname := typeURL
	if slash := strings.LastIndex(mname, "/"); slash >= 0 {
		mname = mname[slash+1:]
	}
	mt := proto.MessageType(mname)
	if mt == nil {
		return nil, fmt.Errorf("unknown message type %q", mname)
	}
	return reflect.New(mt.Elem()).Interface().(proto.Message), nil
}
