package main

import (
	"context"
	"time"

	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/auth"
	cmc "gitlab.com/data-custodian/custodian/components/lib-common/pkg/config"
	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/log"
	clog "gitlab.com/data-custodian/custodian/components/lib-common/pkg/log/context"
	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/signal"

	"github.com/sdsc-ordes/modos-rs/components/kms/internal/config"
	"github.com/sdsc-ordes/modos-rs/components/kms/pkg/service"
	"github.com/sdsc-ordes/modos-rs/components/kms/pkg/storage"
	st "github.com/sdsc-ordes/modos-rs/components/kms/pkg/storage/types"
)

func loadConfigs(configDir string, dataDir string) (conf config.Config) {

	conf, err := cmc.LoadConfigs[config.Config](configDir)
	log.PanicEf(err, "Failed loading config files.")
	conf.WithDataDir(dataDir)

	log.Info("Config file", "config", conf)

	return
}

func main() {
	args := parseArgs()
	conf := loadConfigs(args.ConfigDir, args.DataDir)
	log.Setup(
		log.WithFileLogs(args.LogToFile),
		log.WithForceDevLog(conf.Log.ForceDevLog))

	ctx, stop := signal.WithSignal(clog.Context(context.Background()))
	defer stop()

	clog.Infof(ctx, "Starting KMS server.")

	client, err := storage.NewStorageS3(ctx, &conf.Storage.Connection)
	log.PanicEf(err, "Could not create S3 storage.")

	jwtVerifier, err := createJWTVerifier(ctx, &conf.OIDC)
	log.PanicEf(err, "Could not create JWT verifier.")

	// FIXME: remove.
	c, err := client.NewCredentials(
		ctx,
		[]st.BucketPermission{
			{Path: "bucket-a", Permissions: []st.Permission{st.PermissionRead}},
			{Path: "bucket-b", Permissions: []st.Permission{st.PermissionWrite}},
		},
		1*time.Hour,
	)
	if err != nil {
		log.ErrorE(err, "Credentials could not be created.")
	}
	clog.Info(ctx, "Credentials created.", "creds", c)

	_ = service.Service{Storage: client, JWTVerifier: jwtVerifier}
}

func createJWTVerifier(ctx context.Context, oidcCfg *config.OIDC) (*auth.JWTVerifier, error) {
	const timeout = 30 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return auth.NewJWTVerifier(
		ctx,
		oidcCfg.Issuer,
		oidcCfg.ClientID,
		auth.WithTrustedAlgorithms(oidcCfg.TrustedAlgorithms...),
		auth.WithTrustedAudiences(oidcCfg.TrustedAudiences...),
	)
}
