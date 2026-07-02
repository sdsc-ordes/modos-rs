package main

import (
	"time"

	cmc "gitlab.com/data-custodian/custodian/components/lib-common/pkg/config"
	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/log"

	"gitlab.com/data-custodian/dac-portal/components/nats-server/internal/config"

	"github.com/nats-io/nats-server/v2/server"
)

func loadConfigs(configDir string, dataDir string) (conf config.Config) {
	conf, err := cmc.LoadConfigs[config.Config](configDir)
	log.PanicEf(err, "Failed loading config files.")
	conf.WithDataDir(dataDir)

	log.Info("Config file", "config", conf)

	return
}

func startNATSServer(
	conf *config.Config,
) (ns *server.Server, err error) {
	opts := &server.Options{ //nolint:exhaustruct // Ok.
		Port:      conf.Server.Port,
		JetStream: true,                      // Enable persistence
		StoreDir:  conf.Server.PersistentDir, // Directory for JetStream
	}

	ns, err = server.NewServer(opts)
	if err != nil {
		return nil, err
	}

	ns.Start()

	if !ns.ReadyForConnections(10 * time.Second) { //nolint:mnd
		return nil, err
	}

	return ns, nil
}

func main() {
	args := parseArgs()
	conf := loadConfigs(args.ConfigDir, args.DataDir)
	log.Setup(
		log.WithFileLogs(args.LogToFile),
		log.WithForceDevLog(conf.Log.ForceDevLog))

	ns, err := startNATSServer(&conf)
	log.PanicEf(err, "Could not start NATS server.")
	ns.WaitForShutdown()
}
