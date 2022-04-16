package xcontext

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPairs(t *testing.T) {
	type args struct {
		kv []string
	}
	tests := []struct {
		name           string
		args           args
		want           Metadata
		panicAssertion func(t require.TestingT, f assert.PanicTestFunc, msgAndArgs ...interface{})
	}{
		{
			name: "nil",
			args: args{
				kv: nil,
			},
			want:           Metadata{},
			panicAssertion: require.NotPanics,
		},
		{
			name: "empty",
			args: args{
				kv: []string{},
			},
			want:           Metadata{},
			panicAssertion: require.NotPanics,
		},
		{
			name: "not even",
			args: args{
				kv: []string{"1"},
			},
			panicAssertion: require.Panics,
		},
		{
			name: "even",
			args: args{
				kv: []string{"1", "2"},
			},
			want: Metadata{
				"1": "2",
			},
			panicAssertion: require.NotPanics,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Metadata
			tt.panicAssertion(t, func() {
				got = Pairs(tt.args.kv...)
			})
			require.Equal(t, tt.want, got)
		})
	}
}

func TestJoin(t *testing.T) {
	type args struct {
		mds []Metadata
	}
	tests := []struct {
		name string
		args args
		want Metadata
	}{
		{
			name: "empty",
			args: args{
				mds: []Metadata{},
			},
			want: Metadata{},
		},
		{
			name: "one",
			args: args{
				mds: []Metadata{
					{
						"1": "2",
					},
				},
			},
			want: Metadata{
				"1": "2",
			},
		},
		{
			name: "multi",
			args: args{
				mds: []Metadata{
					{
						"1": "2",
					},
					{
						"1": "3",
						"2": "3",
					},
					{
						"3": "3",
					},
				},
			},
			want: Metadata{
				"1": "3",
				"2": "3",
				"3": "3",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, Join(tt.args.mds...))
		})
	}
}

func TestMetadata_Len(t *testing.T) {
	tests := []struct {
		name string
		m    Metadata
		want int
	}{
		{
			name: "empty",
			m:    Metadata{},
			want: 0,
		},
		{
			name: "one",
			m: Metadata{
				"1": "2",
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.m.Len())
		})
	}
}

func TestMetadata_Get(t *testing.T) {
	type args struct {
		k string
	}
	tests := []struct {
		name string
		m    Metadata
		args args
		want string
	}{
		{
			name: "empty",
			m:    Metadata{},
			args: args{
				k: "1",
			},
			want: "",
		},
		{
			name: "exists",
			m: Metadata{
				"1": "2",
			},
			args: args{
				k: "1",
			},
			want: "2",
		},
		{
			name: "not exists",
			m: Metadata{
				"1": "2",
			},
			args: args{
				k: "3",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.m.Get(tt.args.k))
		})
	}
}

func TestMetadata_Set(t *testing.T) {
	type args struct {
		k   string
		val string
	}
	tests := []struct {
		name string
		m    Metadata
		args args
		want Metadata
	}{
		{
			name: "empty",
			m:    Metadata{},
			args: args{
				k:   "1",
				val: "2",
			},
			want: Metadata{
				"1": "2",
			},
		},
		{
			name: "exists",
			m: Metadata{
				"1": "2",
			},
			args: args{
				k:   "1",
				val: "3",
			},
			want: Metadata{
				"1": "3",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Set(tt.args.k, tt.args.val)
			require.Equal(t, tt.want, tt.m)
		})
	}
}

func TestMetadata_Copy(t *testing.T) {
	tests := []struct {
		name string
		m    Metadata
		want Metadata
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.m.Copy())
		})
	}
}

func TestMetadata_Append(t *testing.T) {
	type args struct {
		k   string
		val string
	}
	tests := []struct {
		name string
		m    Metadata
		args args
		want Metadata
	}{
		{
			name: "empty",
			m:    Metadata{},
			args: args{
				k:   "1",
				val: "2",
			},
			want: Metadata{
				"1": "2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Append(tt.args.k, tt.args.val)
			require.Equal(t, tt.want, tt.m)
		})
	}
}

func TestMetadata_Delete(t *testing.T) {
	type args struct {
		k string
	}
	tests := []struct {
		name string
		md   Metadata
		args args
		want Metadata
	}{
		{
			name: "empty",
			md:   Metadata{},
			args: args{
				k: "1",
			},
			want: Metadata{},
		},
		{
			name: "exists",
			md: Metadata{
				"1": "2",
				"3": "4",
			},
			args: args{
				k: "1",
			},
			want: Metadata{
				"3": "4",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.md.Delete(tt.args.k)
			require.Equal(t, tt.want, tt.md)
		})
	}
}
