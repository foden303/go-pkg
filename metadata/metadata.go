package metadata

import (
	"strings"
)

type Metadata map[string][]string

func New(mds ...map[string][]string) Metadata {
	md := Metadata{}
	for _, m := range mds {
		for k, list := range m {
			md[k] = append(md[k], list...)
		}
	}
	return md
}

func (m Metadata) Add(key, value string) {
	if key == "" {
		return
	}
	keyLower := strings.ToLower(key)
	m[keyLower] = append(m[keyLower], value)
}

func (m Metadata) Get(key string) string {
	keyLower := strings.ToLower(key)
	if values, ok := m[keyLower]; ok && len(values) > 0 {
		return values[0]
	}
	return ""
}

func (m Metadata) Set(key, value string) {
	if key == "" || value == "" {
		return
	}
	m[strings.ToLower(key)] = []string{value}
}

func (m Metadata) Range(f func(key string, value []string) bool) {
	for k, values := range m {
		if !f(k, values) {
			break
		}
	}
}

func (m Metadata) Values(key string) []string {
	return m[strings.ToLower(key)]
}

func (m Metadata) Clone() Metadata {
	md := make(Metadata, len(m))
	for k, v := range m {
		// md[k] = slices.Clone(v)
		md[k] = append([]string{}, v...)
	}
	return md
}
