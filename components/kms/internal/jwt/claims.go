package jwt

import (
	"strings"

	"github.com/sdsc-ordes/modos-rs/components/kms/internal/config"
	"github.com/sdsc-ordes/modos-rs/components/kms/pkg/storage/types"
	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/auth"
	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/errors"
)

type Claims struct {
	*auth.StandardClaims
	BucketPermissions types.BucketPermissions

	cfg config.ClaimBucketPermissions
}

// NewClaims creates a new bucket permissions claim.
func NewClaims(cfg *config.ClaimBucketPermissions) *Claims {
	return &Claims{cfg: *cfg}
}

// InitStdClaims implements [auth.IInitStdClaims] interface.
func (c *Claims) InitStdClaims(stdClaims *auth.StandardClaims) {
	c.StandardClaims = stdClaims
}

// InitCustomClaims implements [auth.IClaimInitializer] interface.
func (c *Claims) InitCustomClaims(getter auth.ClaimGetter) error {
	var bps []any

	err := getter(c.cfg.Name, &bps)
	if err != nil {
		return errors.AddContext(err, "could not convert claim '%v'", c.cfg.Name)
	}

	for _, v := range bps {
		m, ok := v.(map[string]any)
		if !ok {
			return errors.New("could not extract bucket permissions claim")
		}

		p, ok := m[c.cfg.PathName]
		if !ok {
			return errors.New("could not extract bucket permissions claim: '%s'", c.cfg.PathName)
		}
		path, ok := p.(string)
		if !ok {
			return errors.New("bucket permission claim: '%s' not a string", c.cfg.PathName)
		}

		bp, ok := m[c.cfg.PermissionsName]
		if !ok {
			return errors.New(
				"could not extract bucket permissions claim: '%s'",
				c.cfg.PermissionsName,
			)
		}
		permsS, ok := bp.(string)
		if !ok {
			return errors.New("bucket permission claim: '%s' not a string", c.cfg.PermissionsName)
		}

		permsSplit := strings.Split(permsS, ",")
		permissions, err := validatePermissions(permsSplit, &c.cfg)
		if err != nil {
			return err
		}

		c.BucketPermissions = append(c.BucketPermissions, types.BucketPermission{
			Path:        path,
			Permissions: permissions,
		})
	}

	return nil
}

func validatePermissions(
	in []string,
	cfg *config.ClaimBucketPermissions,
) ([]types.Permission, error) {
	var out []types.Permission

	for i := range in {
		switch in[i] {
		case cfg.PermissionsReadTagName:
			out = append(out, types.PermissionRead)
		case cfg.PermissionsWriteTagName:
			out = append(out, types.PermissionWrite)
		default:
			return nil, errors.New(
				"Bucket permissions contains unknown permissions '%v'.",
				in[i],
			)
		}
	}

	return out, nil
}
