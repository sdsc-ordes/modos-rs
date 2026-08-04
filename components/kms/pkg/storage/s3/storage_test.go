//go:build test && unittest

package s3

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
)

func TestS3(t *testing.T) {
	RunSpecs(t, "storage s3")
}
