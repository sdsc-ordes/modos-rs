//go:build test && unittest

package s3

import (
	"encoding/json"

	"github.com/stretchr/testify/require"

	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("credentials", func() {
	var t require.TestingT
	BeforeEach(func() { t = GinkgoT() })

	It("should unmarshal", func() {
		c := S3Credentials{
			AccessKeyID:     "asdf",
			SecretAccessKey: "asdf",
			SessionToken:    "asdf",
		}

		b, err := json.Marshal(c.MarshalerJSON())
		require.NoError(t, err)

		d := S3Credentials{}
		err = json.Unmarshal(b, &d)

		require.NoError(t, err)
		require.Equal(t, "asdf", string(d.AccessKeyID))
	})
})
