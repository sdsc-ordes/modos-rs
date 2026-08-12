package config

import (
	modosregistry "modos-rs/tools/quitsh/pkg/modos-rs/registry"
	"modos-rs/tools/quitsh/pkg/runner/config"
	"path"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	clone "github.com/huandu/go-clone/generic"
	rootcmd "github.com/sdsc-ordes/quitsh/pkg/cli/cmd/root"
	qConfig "github.com/sdsc-ordes/quitsh/pkg/config"
	"github.com/sdsc-ordes/quitsh/pkg/dag"
	fs "github.com/sdsc-ordes/quitsh/pkg/filesystem"
	"github.com/sdsc-ordes/quitsh/pkg/log"
	qRunnerConfig "github.com/sdsc-ordes/quitsh/pkg/runner/config"
	"github.com/sdsc-ordes/quitsh/pkg/toolchain"
)

type (
	CommandArgs struct {
		// Arguments needed to make the root command in `quitsh` work.
		Root rootcmd.Args `yaml:"root"`

		// Arguments needed to make the `execute`
		// command in `quitsh` work. This is used when `quitsh` dispatches over a toolchain
		// and needs to call it self (see `exec.AddCmd`).
		DispatchArgs toolchain.DispatchArgs `yaml:"toolchainDispatch"`

		// Exec arguments.
		ExecArgs dag.ExecArgs `yaml:"execArgs"`
	}

	Config struct {
		// All command arguments of our `quitsh` instance.
		Commands CommandArgs `yaml:"commands"`

		// Some Nix settings.
		Nix qRunnerConfig.NixSettings `yaml:"nix"`

		// The build settings which get copied and injected into the runners:
		// - `custodian::build-go`
		Build config.BuildSettings `yaml:"build"`

		// The lint settings which get copied and injected into the runners:
		// - `custodian::lint-go`
		Lint config.LintSettings `yaml:"lint"`

		// The test settings which get copied and injected into the runners:
		// - `custodian::test-go`
		Test config.TestSettings `yaml:"test"`

		// The images settings which get copied and injected into the runners:
		// - `custodian::image-nix`
		// - `custodian::image-containerfile`
		Image config.ImageSettings `yaml:"image"`

		// The CI settings.
		CI CISettings `yaml:"ci"`
	}

	CISettings struct {
		// If that file is set, it will be used.
		PipelineSettingsFile string `yaml:"pipelineSettingsFile"`
	}
)

// SetDefaults implements [defaults.Setter].
func (c *CISettings) SetDefaults() {
	c.PipelineSettingsFile = DefaultPipelineSettingsFile()
}

func (c *Config) SetDefaults() {
	c.Image.Push.RegistryDomain, c.Image.Push.RegistryBasePathFmt =
		modosregistry.NewRegistryBaseName()

	c.Nix.Cache.SSH.Enable = true
	c.Nix.Cache.SSH.HostName = "nix-cache.swisscustodian.ch"
	c.Nix.Cache.SSH.HostPublicKey = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIMP4u5o2KC6f2OpO4b1GRQzhdcBheLENmS++xw1y5JvY"
	c.Nix.Cache.SSH.Write.PrivateKeyEnv = "NIX_CACHE_SSH_WRITE_KEY"
}

// New returns a Custodian quitsh [config.IConfig] with default values.
func New() (args Config) {
	// Fields which are also flags will be initialized
	// by the flags default values.

	err := defaults.Set(&args)
	log.PanicE(err, "could not default initialize config")

	return
}

func DefaultPipelineSettingsFile() string {
	return path.Join(fs.OutputDir, fs.OutCIDir, "pipeline-settings.yaml")
}

// Clone implements [config.IConfig].
func (c *Config) Clone() qConfig.IConfig {
	return clone.Clone(c)
}

// Validate implements [config.IConfig] interface.
func (c *Config) Validate() error {
	return validator.New().StructExcept(c, "Commands")
}
