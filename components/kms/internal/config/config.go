package config

import "github.com/sdsc-ordes/modos-rs/components/kms/pkg/storage/types"

type Config struct {
	// Some log settings.
	Log Log `yaml:"log"`

	// The server information.
	Server Server `yaml:"server"`

	Storage StorageS3 `yaml:"storage"`
}

type (
	Server struct {
		// The hostname.
		Hostname string `yaml:"hostname" default:"localhost"`

		// The port for the portal endpoints.
		Port int `yaml:"port" default:"3020"`
	}

	StorageS3 struct {
		Connection types.S3Connection `yaml:"connection"`
	}

	Log struct {
		ForceDevLog bool `yaml:"enableDevLog"`
	}
)

// WithDataDir changes the config settings to accommodate
// for the data directory `dir`. Certain configs (paths) are relative to this
// directory if not absolute specified.
func (c *Config) WithDataDir(_ string) {
	// dir = fs.MakeAbsolute(dir)
	// Currently nothing to do here since not relative files.
}
