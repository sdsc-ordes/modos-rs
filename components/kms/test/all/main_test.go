//go:build test && integration

package all

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sdsc-ordes/modos-rs/components/kms/test/service"
	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/log"
)

var stopServices func() error

func TestAll(t *testing.T) {
	log.Setup(log.WithForceDevLog(true))

	RegisterFailHandler(Fail)

	suiteConfig, reporterConfig := GinkgoConfiguration()
	// reporterConfig.FullTrace = true

	RunSpecs(t, "kms", suiteConfig, reporterConfig)
}

var _ = BeforeSuite(func() {
	pcCtx, stop := service.Start()
	stopServices = stop
	service.Wait(pcCtx)
})

var _ = AfterSuite(func() {
	if stopServices != nil {
		Expect(stopServices()).To(Succeed())
	}
})
