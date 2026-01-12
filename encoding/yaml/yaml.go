package yaml

import (
	"go-pkg/encoding"

	"gopkg.in/yaml.v3"
)
// init registers the YAML codec.
func init() {
	encoding.RegisterCodec(codec{})
}

type codec struct{}
// Name returns the name of the codec.
func (codec) Name() string {
	return "yaml"
}
// Marshal marshals v into YAML format.
func (c codec) Marshal(v any) ([]byte, error) {
	return yaml.Marshal(v)
}
// Unmarshal un-marshals data into v from YAML format.
func (c codec) Unmarshal(data []byte, v any) error {
	return yaml.Unmarshal(data, v)
}
