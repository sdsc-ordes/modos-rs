package config

import (
	fs "gitlab.com/data-custodian/custodian/components/lib-common/pkg/filesystem"
)

type Config struct {
	// Some log settings.
	Log Log `yaml:"log"`

	// The server information.
	Server Server `yaml:"server"`
}

type (
	Server struct {
		// The hostname.
		Hostname string `yaml:"hostname" default:"localhost"`

		// The port for the portal endpoints.
		Port int `yaml:"port"     default:"3020"`

		// Where the persistent data of the queue is stored.
		PersistentDir string `yaml:"persistentDir" default:"nats-server"`
	}

	Log struct {
		ForceDevLog bool `yaml:"enableDevLog"`
	}
)

// WithDataDir changes the config settings to accommodate
// for the data directory `dir`. Certain configs (paths) are relative to this
// directory if not absolute specified.
func (c *Config) WithDataDir(dir string) {
	dir = fs.MakeAbsolute(dir)

	c.Server.PersistentDir = fs.MakeAbsoluteTo(dir, c.Server.PersistentDir)
}
