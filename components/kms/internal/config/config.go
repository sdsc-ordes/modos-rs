package config

import "github.com/sdsc-ordes/modos-rs/components/kms/pkg/storage/types"

type Config struct {
	// Some log settings.
	Log Log `yaml:"log"`

	// The server information.
	Server Server `yaml:"server"`

	Storage StorageS3 `yaml:"storage"`

	OIDC OIDC `yaml:"oidc"`
}

type (
	Server struct {
		// The hostname.
		Hostname string `yaml:"hostname" default:"localhost"`

		// The port for the portal endpoints.
		Port int `yaml:"port" default:"3020"`
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
		//`{<PathName>: "...", <PermissionsName>: "..." }`.
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

// WithDataDir changes the config settings to accommodate
// for the data directory `dir`. Certain configs (paths) are relative to this
// directory if not absolute specified.
func (c *Config) WithDataDir(_ string) {
	// dir = fs.MakeAbsolute(dir)
	// Currently nothing to do here since not relative files.
}
