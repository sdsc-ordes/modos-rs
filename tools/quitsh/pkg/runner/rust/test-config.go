package rustrunner

import (
	"github.com/sdsc-ordes/quitsh/pkg/common"
	"github.com/sdsc-ordes/quitsh/pkg/component/step"

	"github.com/creasty/defaults"
)

type RunnerTestConfig struct {
	RunnerConfigBuild

	// Test the specified test targets.
	// If empty all test targets are built.
	Tests []string `yaml:"tests"`

	// Additional arguments forwarded to the test tool (`cargo test <args>`).
	Args []string `yaml:"args"`

	// Additional arguments forwarded to the test executable (`cargo test -- <args>...`).
	TestArgs []string `yaml:"testArgs"`
}

func (c *RunnerTestConfig) Validate() error {
	return common.Validator().Struct(c)
}

// UnmarshalTestConfig unmarshalls for the [RunnerTestConfig].
func UnmarshalTestConfig(raw step.AuxConfigRaw) (step.AuxConfig, error) {
	config := &RunnerTestConfig{} //nolint:exhaustruct // Ok.
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
