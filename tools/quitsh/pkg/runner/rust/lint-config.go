package rustrunner

import (
	"github.com/sdsc-ordes/quitsh/pkg/common"
	"github.com/sdsc-ordes/quitsh/pkg/component/step"

	"github.com/creasty/defaults"
)

type RunnerConfigLint struct {
	RunnerConfigBuild

	// Run clippy per environment type.
	// Defaults to `development`.
	PerEnvs []common.EnvironmentType `yaml:"perEnvs"`
}

func (c *RunnerConfigLint) Validate() error {
	return common.Validator().Struct(c)
}

// UnmarshalLintConfig unmarshalls [RunnerConfigLint].
func UnmarshalLintConfig(raw step.AuxConfigRaw) (step.AuxConfig, error) {
	config := &RunnerConfigLint{} //nolint:exhaustruct // Ok.
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
