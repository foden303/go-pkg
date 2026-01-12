package metadata

import (
	"context"
	"fmt"
)

type clientMetadataKey struct{}

// NewClientContext returns a new context with the provided Metadata attached.
func NewClientContext(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, clientMetadataKey{}, md)
}

// FromClientContext retrieves the Metadata from the context.
func FromClientContext(ctx context.Context) (Metadata, bool) {
	md, ok := ctx.Value(clientMetadataKey{}).(Metadata)
	return md, ok
}

// AppendToClientContext appends key-value pairs to the Metadata in the context.
func AppendToClientContext(ctx context.Context, kv ...string) context.Context {
	if len(kv)%2 == 1 {
		panic(fmt.Sprintf("metadata: AppendToClient got an odd number of input pairs for metadata %d", len(kv)))
	}
	md, _ := FromClientContext(ctx)
	md = md.Clone()
	for i := 0; i < len(kv); i += 2 {
		md.Set(kv[i], kv[i+1])
	}
	return NewClientContext(ctx, md)
}

// MergeToClientContext merges another Metadata into the Metadata in the context.
func MergeToClientContext(ctx context.Context, other Metadata) context.Context {
	md, _ := FromClientContext(ctx)
	md = md.Clone()
	for k, v := range other {
		md[k] = append(md[k], v...)
	}
	return NewClientContext(ctx, md)
}
