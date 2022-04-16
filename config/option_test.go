package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithLocalFile(t *testing.T) {
	type args struct {
		opt LocalOption
	}
	tests := []struct {
		name string
		args args
		want option
	}{
		{
			name: "Default value",
			args: args{
				opt: LocalOption{},
			},
			want: option{
				localEnable: true,
				localDir:    "./etc/conf/",
				localName:   "config",
				localType:   "yaml",
			},
		},
		{
			name: "All custom value",
			args: args{
				opt: LocalOption{
					Directory: "/etc/conf/",
					Filename:  "tmp_config",
					Type:      "json",
				},
			},
			want: option{
				localEnable: true,
				localDir:    "/etc/conf/",
				localName:   "tmp_config",
				localType:   "json",
			},
		},
		{
			name: "Optional custom value",
			args: args{
				opt: LocalOption{
					Filename: "tmp_config",
					Type:     "json",
				},
			},
			want: option{
				localEnable: true,
				localDir:    "./etc/conf/",
				localName:   "tmp_config",
				localType:   "json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got option
			WithLocalFile(tt.args.opt)(&got)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestWithConsul(t *testing.T) {
	type args struct {
		opt ConsulOption
	}
	tests := []struct {
		name string
		args args
		want option
	}{
		{
			name: "Default value",
			args: args{
				opt: ConsulOption{},
			},
			want: option{
				remoteEnable:   true,
				remoteDriver:   "consul",
				remoteEndpoint: "localhost:8500",
				remotePath:     "SERVICE_CONFIG",
				remoteType:     "yaml",
			},
		},
		{
			name: "All custom value",
			args: args{
				opt: ConsulOption{
					Endpoint: "192.168.0.2:8500",
					Path:     "TMP_CONFIG",
					Type:     "json",
				},
			},
			want: option{
				remoteEnable:   true,
				remoteDriver:   "consul",
				remoteEndpoint: "192.168.0.2:8500",
				remotePath:     "TMP_CONFIG",
				remoteType:     "json",
			},
		},
		{
			name: "Optional custom value",
			args: args{
				opt: ConsulOption{
					Path: "TMP_CONFIG",
				},
			},
			want: option{
				remoteEnable:   true,
				remoteDriver:   "consul",
				remoteEndpoint: "localhost:8500",
				remotePath:     "TMP_CONFIG",
				remoteType:     "yaml",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got option
			WithConsul(tt.args.opt)(&got)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestBatchFiles(t *testing.T) {
	type args struct {
		opt BatchFileOption
	}
	tests := []struct {
		name string
		args args
		want option
	}{
		{
			name: "Default value",
			args: args{
				opt: BatchFileOption{},
			},
			want: option{
				localEnable: true,
				localDir:    "./etc/conf/",
				localType:   "yaml",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var o option
			WithBatchFiles(tt.args.opt)(&o)
			require.Equal(t, tt.want, o)
		})
	}
}
