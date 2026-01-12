package json

import (
	"encoding/json"
	"go-pkg/encoding"
	"reflect"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var (
	// MarshalOptions defines the options for marshaling protobuf messages to JSON.
	MarshalOptions = protojson.MarshalOptions{
		EmitUnpopulated: true,
	}
	// UnmarshalOptions defines the options for unmarshaling protobuf messages from JSON.
	UnmarshalOptions = protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}
)

// init registers the JSON codec.
func init() {
	encoding.RegisterCodec(codec{})
}

type codec struct{}

// Name returns the name of the codec.
func (codec) Name() string {
	return "json"
}

// Marshal marshals v into JSON format.
func (c codec) Marshal(v any) ([]byte, error) {
	switch m := v.(type) {
	case proto.Message:
		return MarshalOptions.Marshal(m)
	case json.Marshaler:
		return m.MarshalJSON()
	default:
		return json.Marshal(v)
	}
}

// Unmarshal un-marshals data into v from JSON format.
func (c codec) Unmarshal(data []byte, v any) error {
	switch m := v.(type) {
	case proto.Message:
		return UnmarshalOptions.Unmarshal(data, m)
	case json.Unmarshaler:
		return m.UnmarshalJSON(data)
	default:
		rv := reflect.ValueOf(v)
		for rv := rv; rv.Kind() == reflect.Ptr; {
			if rv.IsNil() {
				rv.Set(reflect.New(rv.Type().Elem()))
			}
			rv = rv.Elem()
		}
		if m, ok := reflect.Indirect(rv).Interface().(proto.Message); ok {
			return UnmarshalOptions.Unmarshal(data, m)
		}
		return json.Unmarshal(data, v)
	}
}
