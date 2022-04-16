package exlog

import (
	"os"
	"testing"

	"github.com/arthurkiller/rollingwriter"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	type args struct {
		o Option
	}
	tests := []struct {
		name string
		args args
		want zerolog.Logger
	}{
		{
			name: "default",
			args: args{
				o: Option{},
			},
			want: zerolog.New(os.Stdout).With().Timestamp().Logger().Level(zerolog.InfoLevel),
		},
		{
			name: "file",
			args: args{
				o: Option{
					File: FileWriter{
						Enable:   true,
						Filename: "tmp",
					},
				},
			},
			want: func() zerolog.Logger {
				conf := rollingwriter.NewDefaultConfig()
				conf.LogPath = "/tmp/logs"
				conf.FileName = "tmp"
				w, _ := rollingwriter.NewWriterFromConfig(&conf)
				logger := zerolog.New(w).With().Timestamp().Logger().Level(zerolog.InfoLevel)
				return logger
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Init(tt.args.o)
			if !tt.args.o.File.Enable {
				require.Equal(t, tt.want, got)
			}
			log.Logger.Debug().Msg("debug")
			log.Logger.Info().Msg("info")
			log.Logger.Warn().Msg("warn")
			if tt.args.o.File.Enable {
				_, err := os.Stat("/tmp/logs/tmp.log")
				require.NoError(t, err)
			}
		})
	}
}
