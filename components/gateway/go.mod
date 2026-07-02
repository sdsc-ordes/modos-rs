module gitlab.com/data-custodian/dac-portal/components/nats-server

go 1.25.0

// If you need a quick way to develop with custodian, uncomment the below line.
// replace gitlab.com/data-custodian/custodian/components/contract-manager => ../../../custodian/components/contract-manager
// replace gitlab.com/data-custodian/custodian/components/lib-common => ../../../custodian/components/lib-common
// replace gitlab.com/data-custodian/custodian/tools/quitsh => ../../../custodian/tools/quitsh

require (
	github.com/alexflint/go-arg v1.5.1
	github.com/creasty/defaults v1.8.0
	github.com/google/uuid v1.6.0
	github.com/hashicorp/go-version v1.7.0
	github.com/nats-io/nats-server/v2 v2.12.5
	github.com/nats-io/nats.go v1.49.0
	gitlab.com/data-custodian/custodian/components/lib-common v0.0.0-20260420063715-d2cf40432642
)

require (
	github.com/alexflint/go-scalar v1.2.0 // indirect
	github.com/antithesishq/antithesis-sdk-go v0.6.0-default-no-op // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/go-tpm v0.9.8 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/klauspost/compress v1.18.4 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/minio/highwayhash v1.0.4-0.20251030100505-070ab1a87a76 // indirect
	github.com/nats-io/jwt/v2 v2.8.0 // indirect
	github.com/nats-io/nkeys v0.4.15 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/sagikazarmark/locafero v0.9.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.14.0 // indirect
	github.com/spf13/cast v1.9.2 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/spf13/viper v1.20.1 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/crypto v0.48.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.34.0 // indirect
	golang.org/x/time v0.15.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
