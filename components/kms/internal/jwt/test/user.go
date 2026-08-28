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
		ID:    uuid.MustParse("9f358a03-c15a-4730-905f-5d6be872cf90"),
		Email: "test@example.com",
	}
}
