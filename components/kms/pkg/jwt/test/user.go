//go:build test && (integration || unittest)

package jwt

import "github.com/google/uuid"

type User struct {
	ID    uuid.UUID
	Email string
}

// DefaultUser returns the default user in the OAuth instances.
func DefaultUser() User {
	return User{ID: uuid.MustParse(""), Email: "test@example.com"}
}
