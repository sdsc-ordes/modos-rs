//go:build test && integration

package all

import (
	"context"
	"testing"

	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/log"
)

type (
	TestContext struct {
		Context context.Context
	}

	TestContextOption func(*testContextOpts)

	testContextOpts struct {
	}
)

// NewTestContext sets up test context.
func NewTestContext(t testing.TB, opts ...TestContextOption) (testCtx *TestContext) {
	log.Info("Construct new test context for contract-manager.")

	var o testContextOpts
	o.Apply(opts...)

	return &TestContext{}
}

func (c *TestContext) Close(t testing.TB) {
}

func (c *testContextOpts) Apply(options ...TestContextOption) {
	for _, f := range options {
		f(c)
	}
}
