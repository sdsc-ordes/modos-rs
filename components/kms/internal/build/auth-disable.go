//go:build auth_disable

package build

import (
	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/build"
)

const (
	DisableAuth = build.DevelopmentEnvEnabled && true
)
