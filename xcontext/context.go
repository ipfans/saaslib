package xcontext

import "context"

var metaKey = struct{}{}

func MetadataFromContext(ctx context.Context) Metadata {
	val := ctx.Value(metaKey)
	if val == nil {
		return Metadata{}
	}
	return val.(Metadata)
}

func AppendContextMetadata(ctx context.Context, m Metadata) context.Context {
	return context.WithValue(ctx, metaKey, Join(MetadataFromContext(ctx), m))
}
