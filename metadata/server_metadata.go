package metadata

import "context"

// serverMetadataKey is the context key for server metadata.
type serverMetadataKey struct{}

// NewServerContext returns a new context with the provided Metadata attached.
func NewServerContext(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, serverMetadataKey{}, md)
}

// FromServerContext retrieves the Metadata from the context.
func FromServerContext(ctx context.Context) (Metadata, bool) {
	md, ok := ctx.Value(serverMetadataKey{}).(Metadata)
	return md, ok
}
