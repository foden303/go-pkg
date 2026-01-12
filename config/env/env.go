package env

type env struct {
	prefixes []string
}

func New(prefixes ...string) *env {
	return &env{prefixes: prefixes}
}
