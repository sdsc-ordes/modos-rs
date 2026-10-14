package config

import (
	"github.com/creasty/defaults"
	"github.com/sdsc-ordes/modos-rs/components/kms/pkg/storage/types"
)

var _ defaults.Setter = (*ServerManagement)(nil)

type Config struct {
	// Some log settings.
	Log Log `yaml:"log"`

	// The server information for the public endpoints.
	Server Server `yaml:"server"`
	// The server information for the management endpoints.
	ServerManagement ServerManagement `yaml:"serverManagement"`

	Storage StorageS3 `yaml:"storage"`

	OIDC OIDC `yaml:"oidc"`
}

type (
	Server struct {
		// The hostname (use 0.0.0.0 for serving it in the cloud).
		Hostname string `yaml:"hostname" default:"localhost"`

		// The port for the public endpoints.
		Port int `yaml:"port" default:"3000"`
	}

	ServerManagement struct {
		Server
	}

	StorageS3 struct {
		Connection types.S3Connection `yaml:"connection"`
	}

	OIDC struct {
		// The OpenID Connect issuer URL.
		Issuer string `yaml:"issuer"`

		// The OpenID Connect client id.
		ClientID string `yaml:"clientID"`

		// Accepted algorithms on the JWT validator.
		TrustedAlgorithms []string `yaml:"trustedAlgorithms" default:"[\"EdDSA\", \"RS256\", \"RS512\", \"ES256\"]"`

		// Accepted audiences (claim: `aud`) on the JWT validator (the `ClientID` is added by default).
		TrustedAudiences []string `yaml:"trustedAudiences"`

		// The bucket permissions claims.
		ClaimBucketPermissions ClaimBucketPermissions `yaml:"claimBucketPermissions"`
	}

	ClaimBucketPermissions struct {
		// The name of the claim with a list
		// of bucket permissions in the form of
		// `{<PathName>: "...", <PermissionsName>: "..." }`.
		Name string `yaml:"name" default:"bps"`

		// The key name of the bucket path.
		PathName string `yaml:"pathName" default:"p"`

		// The key name of the bucket permissions.
		PermissionsName string `yaml:"permissionsName" default:"bp"`

		// The tag name of the read permissions.
		PermissionsReadTagName string `yaml:"permissionsReadTagName" default:"r"`
		// The tag name of the write permissions.
		PermissionsWriteTagName string `yaml:"permissionsWriteTagName" default:"w"`
	}

	Log struct {
		ForceDevLog bool `yaml:"enableDevLog"`
	}
)

// SetDefaults implements [defaults.Setter].
func (s *ServerManagement) SetDefaults() {
	s.Port = 3020
}

// WithDataDir changes the config settings to accommodate
// for the data directory `dir`. Certain configs (paths) are relative to this
// directory if not absolute specified.
func (c *Config) WithDataDir(_ string) {
	// dir = fs.MakeAbsolute(dir)
	// Currently nothing to do here since not relative files.
}
