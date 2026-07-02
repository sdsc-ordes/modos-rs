package rustrunner

import (
	"github.com/sdsc-ordes/quitsh/pkg/common"
	"github.com/sdsc-ordes/quitsh/pkg/component/step"

	"github.com/creasty/defaults"
)

type RunnerConfigBuild struct {
	// Additional build features.
	Features []string `yaml:"features"`

	// Build the specified binary.
	// If empty all binaries are built.
	Binaries []string `yaml:"binaries"`

	// Build the package's library.
	Libraries bool `yaml:"libraries" default:"true"`

	// Build all examples.
	Examples bool `yaml:"examples" default:"true"`

	// Build all tests.
	// If empty all tests are built.
	Tests []string `yaml:"tests"`
}

func (c *RunnerConfigBuild) Validate() error {
	return common.Validator().Struct(c)
}

// UnmarshalBuildConfig unmarshalls [RunnerConfigBuild].
func UnmarshalBuildConfig(raw step.AuxConfigRaw) (step.AuxConfig, error) {
	config := &RunnerConfigBuild{} //nolint:exhaustruct // Ok.
	err := defaults.Set(config)
	if err != nil {
		return nil, err
	}

	// Deserialize if we have something.
	if raw.Unmarshal != nil {
		err = raw.Unmarshal(config)
		if err != nil {
			return nil, err
		}
	}

	err = config.Validate()
	if err != nil {
		return nil, err
	}

	return config, nil
}
