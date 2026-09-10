package s3

import (
	"encoding/json"
	"time"

	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/secret"
)

// S3Credentials represents temp. credentials from an S3.
// It is not serializable.
type S3Credentials struct {
	// The access key ID that identifies the temporary security credentials.
	AccessKeyID secret.RedactedString
	// The secret access key that can be used to sign requests.
	SecretAccessKey secret.RedactedString
	// The token that users must pass to the service API to use the temporary
	// credentials.
	SessionToken secret.RedactedString
	// The date on which the current credentials expire.
	Expiration time.Time
}

type s3credentials struct {
	Type            string    `json:"type"`
	AccessKeyID     string    `json:"accessKeyID"`
	SecretAccessKey string    `json:"secretAccessKey"`
	SessionToken    string    `json:"sessionToken"`
	Expiration      time.Time `json:"expiration"`
}

type marshaller struct {
	*s3credentials
}

// Type implements [types.Credential].
func (c *S3Credentials) Type() string {
	return "s3"
}

// MarshalerJSON implements [types.Credential].
func (c *S3Credentials) MarshalerJSON() json.Marshaler {
	return marshaller{s3credentials: &s3credentials{
		Type:            c.Type(),
		AccessKeyID:     string(c.AccessKeyID),
		SecretAccessKey: string(c.SecretAccessKey),
		SessionToken:    string(c.SessionToken),
	}}
}

func (m marshaller) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.s3credentials)
}

func (c *S3Credentials) MarshalJSON() ([]byte, error) {
	panic("you must not unmarshal S3Credentials")
}
