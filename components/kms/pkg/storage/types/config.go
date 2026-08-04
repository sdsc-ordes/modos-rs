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
		// This is the access key and secret key of a
		// dedicated user who is allowed to call `AssumeRole` to hand out
		// temporary access credentials through the AWS Security Token Service
		// API.
		AccessKey secret.RedactedString `yaml:"accessKey"`
		SecretKey secret.RedactedString `yaml:"secretKey"`

		// The region of the buckets.
		Region string `yaml:"region" default:"us-east-1"`
	}
)
