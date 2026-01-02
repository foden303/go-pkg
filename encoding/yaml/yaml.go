package yaml

import (
	"go-pkg/encoding"

	"gopkg.in/yaml.v3"
)

func init() {
	encoding.RegisterCodec(codec{})
}

type codec struct{}

func (codec) Name() string {
	return "yaml"
}

func (c codec) Marshal(v any) ([]byte, error) {
	return yaml.Marshal(v)
}

func (c codec) Unmarshal(data []byte, v any) error {
	return yaml.Unmarshal(data, v)
}
