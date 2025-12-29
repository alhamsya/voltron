package logging

import (
	"context"
	"sync"
)

type Metadata struct {
	*sync.Map
}

type loggingMetadata struct{}

func (m *Metadata) ToMap() map[string]any {
	res := map[string]any{}

	m.Range(func(key, value any) bool {
		res[key.(string)] = value
		return true
	})

	return res
}

func MetadataFromContext(ctx context.Context) *Metadata {
	val, ok := ctx.Value(loggingMetadata{}).(*Metadata)
	if ok {
		return val
	}

	return &Metadata{&sync.Map{}}
}

func ContextWithMetadata(ctx context.Context, md *Metadata) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if md == nil {
		md = &Metadata{&sync.Map{}}
	}

	return context.WithValue(ctx, loggingMetadata{}, md)
}

func InjectMetadata(ctx context.Context, values map[string]any) context.Context {
	md := MetadataFromContext(ctx)
	for k, v := range values {
		md.Store(k, v)
	}
	return ContextWithMetadata(ctx, md)
}
