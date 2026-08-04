package types

// Credentials represents an agnostic credential.
type Credentials interface {
	Type() string
}
