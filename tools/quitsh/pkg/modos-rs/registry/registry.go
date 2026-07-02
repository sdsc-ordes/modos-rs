package modosregistry

const modosRegistry = "ghcr.io"
const modosImageBasePathFmt = "sdsc-ordes/modos-rs/nix-%s"

// NewRegistryBaseName returns the image base name.
func NewRegistryBaseName() (domain string, baseName string) {
	domain = modosRegistry
	baseName = modosImageBasePathFmt

	return
}
