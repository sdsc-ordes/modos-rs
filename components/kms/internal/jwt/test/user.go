//go:build test && (integration || unittest)

package jwt

import "github.com/google/uuid"

type User struct {
	ID    uuid.UUID
	Email string
}

// DefaultUser returns the default user in the OAuth instances.
func DefaultUser() User {
	return User{
		ID:    uuid.MustParse("81fd14cf-fae8-4d69-9ca1-a3e3b9e03083"),
		Email: "test@example.com",
	}
}
