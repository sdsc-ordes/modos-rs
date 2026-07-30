//go:build test && integration

package service

import (
	"context"
	"path"
	"time"

	"github.com/sdsc-ordes/modos-rs/tools/quitsh/pkg/nix"

	"github.com/sdsc-ordes/quitsh/pkg/ci"
	"github.com/sdsc-ordes/quitsh/pkg/exec/git"
	pc "github.com/sdsc-ordes/quitsh/pkg/exec/process-compose"
	fs "github.com/sdsc-ordes/quitsh/pkg/filesystem"
	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/log"
)

const timeoutServices = 120 * time.Second

// Start starts the services used for the tests.
// If env. variable `MODOS_TEST_SERVICES` is set to `false`
// the services are not started and you can use:
// `just start-services`
// to start the services manually.
func Start() (pcCtx *pc.ProcessComposeCtx, stop func() error) {
	_, rootDir, err := git.NewCtxAtRoot(".")
	log.PanicE(err, "Could not get root directory.")

	pcConfig := "test-services"

	socketPathFile := path.Join(rootDir,
		fs.OutputDir, fs.OutRunDir, "process-compose",
		"pc-"+pcConfig+".sock")

	mustBeStarted := false
	if fs.Exists(socketPathFile) {
		log.Infof("Skip starting the test services, socket '%v' exists.", socketPathFile)
		mustBeStarted = true
	} else {
		log.Info("Starting test services.")
	}

	pcCtx, err = pc.Start(
		log.GetLogger(),
		rootDir,
		nix.DefaultFlakeDirRel,
		pcConfig,
		pc.ProcessComposeOverServicesFlake,
		pc.WithMustBeStarted(mustBeStarted),
		pc.WithSocketPathFile(socketPathFile),
	)

	log.PanicE(err, "Could not start process-compose.")

	if !mustBeStarted {
		// When we started it, we also gonna destroy it.
		stop = func() error {
			log.Info("Teardown test services.")

			return pcCtx.Stop()
		}
	}

	return pcCtx, stop
}

func Wait(pcCtx *pc.ProcessComposeCtx) {
	timeout := timeoutServices
	if ci.IsRunning() {
		timeout = 2 * timeoutServices //nolint:mnd
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	log.Info("Waiting for services to be ready.")
	condFulfilled, err := pcCtx.WaitTill(
		ctx,
		log.GetLogger(),
		pc.ProcessCond{Name: "keycloak", State: pc.ProcessReady},
		pc.ProcessCond{Name: "authentik", State: pc.ProcessReady},
		pc.ProcessCond{Name: "mailhog", State: pc.ProcessReady},
		pc.ProcessCond{Name: "rustfs", State: pc.ProcessReady},
	)

	if err != nil {
		log.PanicE(err, "Failed to wait for 'rustfs', 'keycloak', 'mailhog', 'authentik'.")
	} else if !condFulfilled {
		log.Panic("Failed to fulfill wait condition. Timedout.")
	}
}
