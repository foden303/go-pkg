package encoding

import "strings"

// Codec defines the interface for encoding and decoding data.
type Codec interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
	Name() string
}

var registeredCodecs = make(map[string]Codec)

// RegisterCodec registers a codec with the given name.
func RegisterCodec(codec Codec) {
	if codec == nil {
		panic("cannot register a nil Codec")
	}
	if codec.Name() == "" {
		panic("cannot register Codec with empty string result for Name()")
	}
	registeredCodecs[strings.ToLower(codec.Name())] = codec
}

// GetCodec retrieves a registered codec by name.
func GetCodec(name string) Codec {
	return registeredCodecs[name]
}
