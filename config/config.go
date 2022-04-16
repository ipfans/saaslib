package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ipfans/saaslib/liberrors"
	"github.com/morikuni/failure"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

type Config interface {
	Unmarshal(interface{}) error
}

type config struct {
	local  *viper.Viper
	remote *viper.Viper

	v *viper.Viper
}

// Init returns a new config instance.
// Default args:
//   - ./etc/conf/config.yaml
// nolint:nakedret
func Init(opt ...Option) (conf Config, err error) {
	var o option
	c := &config{
		local:  viper.New(),
		remote: viper.New(),

		v: viper.New(),
	}

	if len(opt) == 0 {
		opt = append(opt, WithLocalFile(LocalOption{}))
	}
	for i := range opt {
		opt[i](&o)
	}

	defer func() {
		if err != nil {
			err = failure.Wrap(err, failure.WithCode(liberrors.ErrConfigReadFailed))
		}
	}()

	if o.localEnable {
		if err = c.readLocalConfig(o); err != nil {
			return
		}
	}

	if o.remoteEnable {
		if err = c.readRemoteConfig(o); err != nil {
			return
		}
	}

	err = c.mergeConfig(o)

	conf = c
	return
}

func (c *config) readSingleFile(o option) (err error) {
	if err = c.local.ReadInConfig(); err != nil {
		err = failure.Wrap(err,
			failure.Context{"config": filepath.Join(o.localDir, o.localName+"."+o.localType)},
		)
		return
	}
	return
}

//nolint:nilerr
func (c *config) batchFiles(o option) error {
	err := filepath.Walk(o.localDir, func(path string, info os.FileInfo, err error) (e error) {
		if info.IsDir() || err != nil {
			return
		}
		fn := strings.Split(filepath.Base(path), ".")[0]
		if fn == "" {
			return
		}

		v := viper.New()
		v.AddConfigPath(o.localDir)
		v.SetConfigName(fn)
		v.SetConfigType(o.localType)
		if e = v.ReadInConfig(); e != nil {
			e = failure.Wrap(e,
				failure.Context{
					"config": filepath.Join(o.localDir, fn+"."+o.localType),
				},
			)
			return
		}
		e = c.v.MergeConfigMap(v.AllSettings())
		e = failure.Wrap(e,
			failure.Message("Merge config failed"),
			failure.Context{
				"config": filepath.Join(o.localDir, fn+"."+o.localType),
			},
		)
		return
	})
	return err
}

func (c *config) readLocalConfig(o option) (err error) {
	c.local = viper.New()
	c.local.AddConfigPath(o.localDir)
	c.local.SetConfigType(o.localType)

	if o.localName != "" {
		c.local.SetConfigName(o.localName)
		err = c.readSingleFile(o)
		return
	} else {
		err = c.batchFiles(o)
	}

	return
}

func (c *config) readRemoteConfig(o option) (err error) {
	if err = c.remote.AddRemoteProvider(o.remoteDriver, o.remoteEndpoint, o.remotePath); err != nil {
		err = failure.Wrap(err,
			failure.Context{
				"driver":   o.remoteDriver,
				"endpoint": o.remoteEndpoint,
				"path":     o.remotePath,
			},
		)
		return
	}
	c.remote.SetConfigType(o.remoteType)
	if err = c.remote.ReadRemoteConfig(); err != nil {
		err = failure.Wrap(err,
			failure.Context{
				"driver":   o.remoteDriver,
				"endpoint": o.remoteEndpoint,
				"path":     o.remotePath,
			},
		)
		return
	}
	return
}

func (c *config) mergeConfig(o option) (err error) {
	if o.localEnable {
		if err = c.v.MergeConfigMap(c.local.AllSettings()); err != nil {
			err = failure.Wrap(err,
				failure.Context{
					"config": filepath.Join(o.localDir, o.localName+"."+o.localType),
				},
			)
			return
		}
	}
	if o.remoteEnable {
		if err = c.v.MergeConfigMap(c.remote.AllSettings()); err != nil {
			err = failure.Wrap(err,
				failure.Context{
					"driver":   o.remoteDriver,
					"endpoint": o.remoteEndpoint,
					"path":     o.remotePath,
				},
			)
			return
		}
	}
	return
}

func (c *config) Unmarshal(v interface{}) error {
	return c.v.Unmarshal(v)
}
