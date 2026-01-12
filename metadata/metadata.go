package metadata

import (
	"strings"
)

// Metadata represents a collection of key-value pairs where each key can have multiple values.
type Metadata map[string][]string

// New creates a new Metadata instance, optionally merging multiple maps into one.
func New(mds ...map[string][]string) Metadata {
	md := Metadata{}
	for _, m := range mds {
		for k, list := range m {
			md[k] = append(md[k], list...)
		}
	}
	return md
}

// Add appends a value to the list of values for the given key.
func (m Metadata) Add(key, value string) {
	if key == "" {
		return
	}
	keyLower := strings.ToLower(key)
	m[keyLower] = append(m[keyLower], value)
}

// Get retrieves the first value associated with the given key.
func (m Metadata) Get(key string) string {
	keyLower := strings.ToLower(key)
	if values, ok := m[keyLower]; ok && len(values) > 0 {
		return values[0]
	}
	return ""
}

// Set sets the value for the given key, replacing any existing values.
func (m Metadata) Set(key, value string) {
	if key == "" || value == "" {
		return
	}
	m[strings.ToLower(key)] = []string{value}
}

// Range iterates over all key-value pairs in the Metadata, calling the provided function for each pair.
func (m Metadata) Range(f func(key string, value []string) bool) {
	for k, values := range m {
		if !f(k, values) {
			break
		}
	}
}

// Values retrieves all values associated with the given key.
func (m Metadata) Values(key string) []string {
	return m[strings.ToLower(key)]
}

// Clone creates a deep copy of the Metadata.
func (m Metadata) Clone() Metadata {
	md := make(Metadata, len(m))
	for k, v := range m {
		// md[k] = slices.Clone(v)
		md[k] = append([]string{}, v...)
	}
	return md
}
