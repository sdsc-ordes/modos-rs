package config

import (
	modosregistry "modos-rs/tools/quitsh/pkg/modos-rs/registry"

	"github.com/creasty/defaults"
	"github.com/sdsc-ordes/quitsh/pkg/log"
	cnConfig "gitlab.com/data-custodian/custodian/tools/quitsh/pkg/config"
)

// Config is a wrapped (newtype idiom) with some addaptions for our use case.
type Config struct {
	cnConfig.Config `yaml:",inline"`
}

// New returns a new config.
func New() (args Config) {
	err := defaults.Set(&args)
	log.PanicE(err, "could not default initialize config")

	return
}

// SetDefaults implements [defaults.Setter] interface.
func (s *Config) SetDefaults() {
	s.Image.Push.RegistryDomain, s.Image.Push.RegistryBasePathFmt =
		modosregistry.NewRegistryBaseName()
}
