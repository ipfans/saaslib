package config

type Option func(o *option)

type option struct {
	localEnable bool
	localDir    string
	localName   string
	localType   string

	remoteEnable   bool
	remoteDriver   string
	remoteEndpoint string
	remotePath     string
	remoteType     string
}

type LocalOption struct {
	Directory string // Directory of the local file. Default: "./etc/conf/"
	Filename  string // Filename without ext. Default: "config"
	Type      string // File type of the local file(yaml/toml/json). Default: "yaml"
}

// WithLocalFile sets the local file path.
// dir: the directory of the local file.
// name: the name of the local file without ext.
// ftype: the file type of the local file. (yaml/toml/json)
func WithLocalFile(opt LocalOption) Option {
	return func(o *option) {
		o.localEnable = true
		if opt.Directory == "" {
			opt.Directory = "./etc/conf/"
		}
		if opt.Filename == "" {
			opt.Filename = "config"
		}
		if opt.Type == "" {
			opt.Type = "yaml"
		}
		o.localDir = opt.Directory
		o.localName = opt.Filename
		o.localType = opt.Type
	}
}

// BatchedFilesOption is the option of the batch files.
type BatchFileOption struct {
	Directory string // Directory of the local file. Default: "./etc/conf/"
	Type      string // File type of the local file(yaml/toml/json). Default: "yaml"
}

// WithBatchFiles sets the batch files. Using lexical order to load the files and merge them.
func WithBatchFiles(opt BatchFileOption) Option {
	return func(o *option) {
		if opt.Directory == "" {
			opt.Directory = "./etc/conf/"
		}
		if opt.Type == "" {
			opt.Type = "yaml"
		}
		o.localEnable = true
		o.localDir = opt.Directory
		o.localType = opt.Type
	}
}

type ConsulOption struct {
	Endpoint string // the consul endpoint url. Default: "localhost:8500"
	Path     string // the consul key. Default: "SERVICE_CONFIG"
	Type     string // the file type of the remote config(yaml/toml/json). Default: "yaml"
}

// WithConsul sets the consul remote config.
// url: the consul url.
// ftype: the file type of the remote config. (yaml/toml/json)
func WithConsul(opt ConsulOption) Option {
	return func(o *option) {
		if opt.Endpoint == "" {
			opt.Endpoint = "localhost:8500"
		}
		if opt.Path == "" {
			opt.Path = "SERVICE_CONFIG"
		}
		if opt.Type == "" {
			opt.Type = "yaml"
		}
		o.remoteEnable = true
		o.remoteDriver = "consul"
		o.remoteEndpoint = opt.Endpoint
		o.remotePath = opt.Path
		o.remoteType = opt.Type
	}
}
