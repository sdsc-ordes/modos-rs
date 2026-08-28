package jwt

import (
	"strings"

	"github.com/sdsc-ordes/modos-rs/components/kms/internal/config"
	"github.com/sdsc-ordes/modos-rs/components/kms/pkg/storage/types"
	"github.com/sdsc-ordes/quitsh/pkg/errors"
	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/auth"
)

type Claims struct {
	*auth.StandardClaims
	BucketPermission types.BucketPermissions

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
	var bp []any

	err := getter(c.cfg.Name, &bp)
	if err != nil {
		return errors.AddContext(err, "could not convert claim '%v'", c.cfg.Name)
	}

	for _, v := range bp {
		m, ok := v.(map[string]any)
		if !ok {
			return errors.New("could not extract bucket permissions claim")
		}

		p, ok := m[c.cfg.PathName]
		if !ok {
			return errors.New("could not extract bucket permissions claim: 'p'")
		}
		path, ok := p.(string)
		if !ok {
			return errors.New("bucket permission claim: 'p' not a string")
		}

		bp, ok := m[c.cfg.PermissionsName]
		if !ok {
			return errors.New("could not extract bucket permissions claim: 'bp'")
		}
		permsS, ok := bp.(string)
		if !ok {
			return errors.New("bucket permission claim: 'bp' not a string")
		}

		permsSplit := strings.Split(permsS, ",")
		permissions, err := validatePermissions(permsSplit)
		if err != nil {
			return err
		}

		c.BucketPermission = append(c.BucketPermission, types.BucketPermission{
			Path:        path,
			Permissions: permissions,
		})
	}

	return nil
}

func validatePermissions(in []string) (out []types.Permission, _ error) {
	for i := range in {
		switch in[i] {
		case "r":
			out = append(out, types.PermissionRead)
		case "w":
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
