package exlog

import (
	"io"
	"os"
	"strings"

	"github.com/arthurkiller/rollingwriter"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type FileWriter struct {
	Enable    bool   // Enable file writer. It will override stdout as default writer.
	Directory string // Output file path. Default is `/tmp/logs/`
	Filename  string // Output file name, it will append ext automaticlly. Default is `app`

	DisableCompress    bool   // Disable compress the trucked file.
	RollingTimePattern string // Rolling time pattern. Visit https://crontab.guru/ for more info. Default `0 0 * * *`
}

type Option struct {
	Level string // debug, info, warn, error, fatal. default: info
	File  FileWriter
}

var (
	defaultLevel = "info"
	logLevel     = map[string]zerolog.Level{
		"debug": zerolog.DebugLevel,
		"info":  zerolog.InfoLevel,
		"warn":  zerolog.WarnLevel,
		"error": zerolog.ErrorLevel,
		"fatal": zerolog.FatalLevel,
	}
)

func rollingWriter(o Option) io.Writer {
	if o.File.Directory == "" {
		o.File.Directory = "/tmp/logs"
	}
	if o.File.Filename == "" {
		o.File.Filename = "app"
	}
	if o.File.RollingTimePattern == "" {
		o.File.RollingTimePattern = "0 0 * * *"
	}
	conf := rollingwriter.NewDefaultConfig()
	conf.LogPath = o.File.Directory
	conf.FileName = o.File.Filename
	conf.RollingPolicy = rollingwriter.TimeRolling
	conf.RollingTimePattern = o.File.RollingTimePattern
	conf.Compress = !o.File.DisableCompress
	wr, err := rollingwriter.NewWriterFromConfig(&conf)
	if err != nil {
		panic(err)
	}
	return wr
}

// Init zerolog.Logger.
// Notice: it will override the default `zerolog/log.Logger`.
func Init(o Option) (logger zerolog.Logger) {
	level, ok := logLevel[strings.ToLower(o.Level)]
	if !ok {
		level = logLevel[defaultLevel]
	}

	var w io.Writer = os.Stdout
	if o.File.Enable {
		w = rollingWriter(o)
	}

	logger = zerolog.New(w).With().Timestamp().Logger().Level(level)
	log.Logger = logger

	return
}
