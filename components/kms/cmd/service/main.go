package main

import (
	"time"

	cmc "gitlab.com/data-custodian/custodian/components/lib-common/pkg/config"
	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/log"

	"github.com/sdsc-ordes/modos-rs/components/kms/internal/config"
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

	log.Infof("Starting KMS server.")
	time.Sleep(100000 * time.Second) //nolint:mnd
}
