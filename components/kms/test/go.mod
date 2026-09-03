module github.com/sdsc-ordes/modos-rs/components/kms/test

go 1.26

// If you need a quick way to develop with custodian, uncomment the below line.
// replace gitlab.com/data-custodian/custodian/components/lib-common => ../../../custodian/components/lib-common

// replace github.com/sdsc-ordes/quitsh => ../../../../quitsh

replace github.com/sdsc-ordes/modos-rs/tools/quitsh => ../../../tools/quitsh

require (
	github.com/onsi/ginkgo/v2 v2.32.0
	github.com/onsi/gomega v1.40.0
	github.com/sdsc-ordes/modos-rs/tools/quitsh v0.0.0-00010101000000-000000000000
	github.com/sdsc-ordes/quitsh v0.44.0
	gitlab.com/data-custodian/custodian/components/lib-common v0.0.0-20260727115656-de5a9ee9d6c1
)

require (
	charm.land/lipgloss/v2 v2.0.5 // indirect
	charm.land/log/v2 v2.0.1 // indirect
	deedles.dev/xiter v0.2.1 // indirect
	github.com/Masterminds/semver/v3 v3.4.0 // indirect
	github.com/bmatcuk/doublestar/v4 v4.8.1 // indirect
	github.com/charlievieth/fastwalk v1.0.10 // indirect
	github.com/charmbracelet/colorprofile v0.4.3 // indirect
	github.com/charmbracelet/ultraviolet v0.0.0-20251205161215-1948445e3318 // indirect
	github.com/charmbracelet/x/ansi v0.11.8 // indirect
	github.com/charmbracelet/x/term v0.2.2 // indirect
	github.com/charmbracelet/x/termios v0.1.1 // indirect
	github.com/charmbracelet/x/windows v0.2.2 // indirect
	github.com/clipperhouse/displaywidth v0.11.0 // indirect
	github.com/clipperhouse/uax29/v2 v2.7.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.13 // indirect
	github.com/go-logfmt/logfmt v0.6.1 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.30.2 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/pprof v0.0.0-20260402051712-545e8a4df936 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-version v1.7.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.4.1 // indirect
	github.com/mattn/go-isatty v0.0.24 // indirect
	github.com/mattn/go-runewidth v0.0.29 // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/otiai10/copy v1.14.1 // indirect
	github.com/otiai10/mint v1.6.3 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/xo/terminfo v1.0.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/crypto v0.55.0 // indirect
	golang.org/x/exp v0.0.0-20260824195058-e88cd73687aa // indirect
	golang.org/x/mod v0.39.0 // indirect
	golang.org/x/net v0.58.0 // indirect
	golang.org/x/sync v0.22.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.41.0 // indirect
	golang.org/x/tools v0.49.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
)
