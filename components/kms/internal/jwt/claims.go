package jwt

import (
	"strings"

	"github.com/sdsc-ordes/modos-rs/components/kms/pkg/storage/types"
	"github.com/sdsc-ordes/quitsh/pkg/errors"
	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/auth"
)

const BucketPermissionsName = "bp"

type Claims struct {
	*auth.StandardClaims
	BucketPermission types.BucketPermissions
}

// InitStdClaims implements [auth.IInitStdClaims] interface.
func (c *Claims) InitStdClaims(stdClaims *auth.StandardClaims) {
	c.StandardClaims = stdClaims
}

// InitCustomClaims implements [auth.IClaimInitializer] interface.
func (c *Claims) InitCustomClaims(getter auth.ClaimGetter) error {
	var bp []any

	err := getter("bp", &bp)
	if err != nil {
		return errors.AddContext(err, "could not convert subject to UUID for '%v'", c.Subject)
	}

	for _, v := range bp {
		m, ok := v.(map[string]any)
		if !ok {
			return errors.New("could not extract bucket permissions claim")
		}

		p, ok := m["p"]
		if !ok {
			return errors.New("could not extract bucket permissions claim: 'p'")
		}
		path, ok := p.(string)
		if !ok {
			return errors.New("bucket permission claim: 'p' not a string")
		}

		bp, ok := m["bp"]
		if !ok {
			return errors.New("could not extract bucket permissions claim: 'bp'")
		}
		permsS, ok := bp.(string)
		if !ok {
			return errors.New("bucket permission claim: 'bp' not a string")
		}
		perms := strings.Split(permsS, ",")

		c.BucketPermission = append(c.BucketPermission, types.BucketPermission{
			Path:        path,
			Permissions: perms,
		})
	}

	return nil
}
