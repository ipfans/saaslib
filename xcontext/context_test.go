package xcontext

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMetadataContext(t *testing.T) {
	type args struct {
		preMD Metadata
		fixMD Metadata
	}
	tests := []struct {
		name string
		args args
		want Metadata
	}{
		{
			name: "empty",
			args: args{
				preMD: nil,
				fixMD: nil,
			},
			want: Metadata{},
		},
		{
			name: "pre",
			args: args{
				preMD: Metadata{
					"a": "a",
				},
				fixMD: nil,
			},
			want: Metadata{
				"a": "a",
			},
		},
		{
			name: "fix",
			args: args{
				preMD: nil,
				fixMD: Metadata{
					"b": "b",
				},
			},
			want: Metadata{
				"b": "b",
			},
		},
		{
			name: "pre and fix",
			args: args{
				preMD: Metadata{
					"a": "a",
				},
				fixMD: Metadata{
					"b": "b",
				},
			},
			want: Metadata{
				"a": "a",
				"b": "b",
			},
		},
		{
			name: "pre and fix with override",
			args: args{
				preMD: Metadata{
					"a": "a",
				},
				fixMD: Metadata{
					"a": "b",
					"b": "b",
				},
			},
			want: Metadata{
				"a": "b",
				"b": "b",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.args.preMD != nil {
				ctx = context.WithValue(ctx, metaKey, tt.args.preMD)
			}
			if tt.args.fixMD != nil {
				ctx = AppendContextMetadata(ctx, tt.args.fixMD)
			}
			require.Equal(t, tt.want, MetadataFromContext(ctx))
		})
	}
}
