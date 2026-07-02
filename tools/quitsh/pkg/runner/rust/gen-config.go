package rustrunner

import (
	"github.com/sdsc-ordes/quitsh/pkg/common"
	"github.com/sdsc-ordes/quitsh/pkg/component/step"

	"github.com/creasty/defaults"
)

type RunnerConfigBuildGen struct {
	// Additional build features.
	Features []string `yaml:"features"`

	// Run the specified binary.
	Binaries []string `yaml:"binaries"`
}

func (c *RunnerConfigBuildGen) Validate() error {
	return common.Validator().Struct(c)
}

// UnmarshalBuildGenConfig unmarshalls [RunnerConfigBuildGen].
func UnmarshalBuildGenConfig(raw step.AuxConfigRaw) (step.AuxConfig, error) {
	config := &RunnerConfigBuildGen{} //nolint:exhaustruct // Ok.
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
