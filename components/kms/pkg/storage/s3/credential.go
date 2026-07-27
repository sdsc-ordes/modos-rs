package s3

import (
	"time"

	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/secret"
)

// S3Credentials represents temp. credentials from an S3.
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

func (c *S3Credentials) Type() string {
	return "s3"
}
