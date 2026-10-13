package types

import "encoding/json"

// Credentials represents an agnostic credential.
// This interface is by design not serializable and only for in-ram use.
// Must be converted to an `json.Marshaller` with
// [Credentials.MarshalerJSON] to be able to serialize it.
type Credentials interface {
	Type() string
	MarshalerJSON() json.Marshaler
}
