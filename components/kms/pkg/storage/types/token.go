package types

import (
	"time"

	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/secret"
)

// Credentials represents temp. credentials.
type Credentials struct {
	// The access key ID that identifies the temporary security credentials.
	AccessKeyID secret.RedactedString
	// The secret access key that can be used to sign requests.
	SecretAccessKey secret.RedactedString

	// The date on which the current credentials expire.
	Expiration time.Time

	// The token that users must pass to the service API to use the temporary
	// credentials.
	SessionToken secret.RedactedString
}
