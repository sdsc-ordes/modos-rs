package types

import (
	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/net"
	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/secret"
)

type (
	S3Connection struct {
		// The endpoint URL.
		Endpoint net.URL `yaml:"endpoint"`
		// For local S3 (garage) expect path-style addressing.
		UsePathStyle bool `yaml:"usePathStyle" default:"false"`

		// Credentials
		AccessKey secret.RedactedString `yaml:"accessKey"`
		SecretKey secret.RedactedString `yaml:"secretKey"`

		// The region of the buckets.
		Region string `yaml:"region" default:"us-east-1"`
	}
)
