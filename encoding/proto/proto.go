package proto

import (
	"errors"
	"go-pkg/encoding"
	"reflect"

	"google.golang.org/protobuf/proto"
)

func init() {
	encoding.RegisterCodec(codec{})
}

type codec struct{}

// Name returns the name of the codec.
func (codec) Name() string {
	return "proto"
}

// Marshal marshals v into ProtoBuf format.
func (c codec) Marshal(v any) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}

// Unmarshal un-marshals data into v from ProtoBuf format.
func (c codec) Unmarshal(data []byte, v any) error {
	pm, err := getPm(v)
	if err != nil {
		return err
	}
	return proto.Unmarshal(data, pm)
}

// getPm retrieves the proto.Message from the given value.
func getPm(v any) (proto.Message, error) {
	if msg, ok := v.(proto.Message); ok {
		return msg, nil
	}
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return nil, errors.New("not proto message")
	}
	val = val.Elem()
	return getPm(val.Interface())
}
