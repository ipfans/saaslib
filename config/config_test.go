package config

import (
	"os"
	"testing"

	"github.com/hashicorp/consul/sdk/testutil"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	server, err := testutil.NewTestServerConfigT(t, nil)
	if err != nil {
		t.Skip("Skip tests because consul server is not available")
	}
	if server.Config.Bootstrap {
		server.WaitForLeader(t)
	}

	defer func() {
		_ = server.Stop()
	}()

	b, err := os.ReadFile("./testfixtures/remote.yaml")
	require.NoError(t, err)
	server.SetKV(t, "TESTCONFIG", b)

	type args struct {
		wantRemote bool
		opt        []Option
	}
	type configKV struct {
		key   string
		value string
	}
	tests := []struct {
		name     string
		args     args
		wantConf []configKV

		assertion require.ErrorAssertionFunc
	}{
		{
			name: "Only local config",
			args: args{
				opt: []Option{
					WithLocalFile(LocalOption{
						Directory: "./testfixtures/",
						Filename:  "local",
					}),
				},
			},
			wantConf: []configKV{
				{
					"a.b.c",
					"1",
				},
			},
			assertion: require.NoError,
		},
		{
			name: "Local config not found",
			args: args{
				opt: []Option{
					WithLocalFile(LocalOption{
						Directory: "./testfixtures/",
						Filename:  "not_exists",
					}),
				},
			},
			assertion: require.Error,
		},
		{
			name: "Only remote config",
			args: args{
				wantRemote: true,
				opt: []Option{
					WithConsul(ConsulOption{
						Endpoint: server.HTTPAddr,
						Path:     "TESTCONFIG",
						Type:     "yaml",
					}),
				},
			},
			wantConf: []configKV{
				{
					"a.b.c",
					"2",
				},
			},
			assertion: require.NoError,
		},
		{
			name: "Only invalid endpoint remote config",
			args: args{
				wantRemote: true,
				opt: []Option{
					WithConsul(ConsulOption{
						Endpoint: "1234",
						Path:     "TESTCONFIG",
						Type:     "yaml",
					}),
				},
			},
			assertion: require.Error,
		},
		{
			name: "Batch files",
			args: args{
				opt: []Option{
					WithBatchFiles(BatchFileOption{
						Directory: "./testfixtures/",
						Type:      "yaml",
					}),
				},
			},
			wantConf: []configKV{
				{
					"a.b.c",
					"2",
				},
				{
					"a.b.d",
					"2",
				},
			},
			assertion: require.NoError,
		},
		{
			name: "Local config and remote config",
			args: args{
				wantRemote: true,
				opt: []Option{
					WithLocalFile(LocalOption{
						Directory: "./testfixtures/",
						Filename:  "local",
					}),
					WithConsul(ConsulOption{
						Endpoint: server.HTTPAddr,
						Path:     "TESTCONFIG",
						Type:     "yaml",
					}),
				},
			},
			wantConf: []configKV{
				{
					"a.b.d",
					"2",
				},
				{
					"a.b.c",
					"2",
				},
			},
			assertion: require.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotConf, err := Init(tt.args.opt...)
			tt.assertion(t, err)
			for i := range tt.wantConf {
				require.Equal(t, tt.wantConf[i].value, gotConf.(*config).v.GetString(tt.wantConf[i].key))
			}
		})
	}
}

func Test_Unmarshal(t *testing.T) {
	type TestConfig struct {
		A struct {
			B struct {
				C string
				D string
			}
		}
	}

	server, err := testutil.NewTestServerConfigT(t, nil)
	if err != nil {
		t.Skip("Skip tests because consul server is not available")
	}
	if server.Config.Bootstrap {
		server.WaitForLeader(t)
	}

	defer func() {
		_ = server.Stop()
	}()

	b, err := os.ReadFile("./testfixtures/remote.yaml")
	require.NoError(t, err)
	server.SetKV(t, "TESTCONFIG", b)
	conf, err := Init(
		WithConsul(ConsulOption{
			Endpoint: server.HTTPAddr,
			Path:     "TESTCONFIG",
			Type:     "yaml",
		}),
		WithLocalFile(LocalOption{
			Directory: "./testfixtures/",
			Filename:  "local",
		}),
	)
	require.NoError(t, err)
	var c TestConfig
	require.NoError(t, conf.Unmarshal(&c))
	require.Equal(t, "2", c.A.B.C)
	require.Equal(t, "2", c.A.B.D)
}
