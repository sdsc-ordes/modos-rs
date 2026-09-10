//go:build test && (integration || unittest)

// Package jwt contains JWT testing functions and constructors.
package jwt

import (
	"log"
	"strings"

	"github.com/lestrrat-go/jwx/v3/jwt"
	"github.com/sdsc-ordes/modos-rs/components/kms/internal/config"
	"github.com/sdsc-ordes/modos-rs/components/kms/pkg/storage/types"
	st "github.com/sdsc-ordes/modos-rs/components/kms/pkg/storage/types"
)

// WithBucketPermissionsClaim adds the bucket permissions claims by the specification in `cfg`
// from given `permissions`.
func WithBucketPermissionsClaim(
	cfg *config.ClaimBucketPermissions,
	permissions st.BucketPermissions,
) Option {
	claim := []any{}

	for _, p := range permissions {
		var perms []string

		for _, permTag := range p.Permissions {
			switch permTag {
			case types.PermissionRead:
				perms = append(perms, cfg.PermissionsReadTagName)
			case types.PermissionWrite:
				perms = append(perms, cfg.PermissionsWriteTagName)
			default:
				log.Panicf("Permissions tag '%v' not implemented.", permTag)
			}
		}

		claim = append(claim,
			map[string]any{
				cfg.PathName:        p.Path,
				cfg.PermissionsName: strings.Join(perms, ","),
			},
		)
	}

	mod := func(b *jwt.Builder) {
		b.Claim(cfg.Name, claim)
	}

	return WithModifications(mod)
}
