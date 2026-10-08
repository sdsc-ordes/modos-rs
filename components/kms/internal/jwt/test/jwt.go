//go:build test && (integration || unittest)

// Package jwt contains JWT testing functions and constructors.
package jwt

import (
	"time"

	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"github.com/stretchr/testify/require"
)

const DefaultIssuerKeycloak = "http://localhost:8081/realms/modos"

const DefaultIssuerAuthentik = "http://localhost:9001/application/o/modos-cli-app"

// DefaultAudience is the same as Keycloak client or the application `modos-cli` in Authentik.
const DefaultAudience = "modos-cli"
const DefaultScope = "openid permissions"

type (
	CustomClaim struct {
		Key   string
		Value any
	}

	Option func(*opts)

	opts struct {
		// Add modifications.
		mods func(b *jwt.Builder)

		// Make the valid -> expired.
		expired bool
	}
)

func NewTokenBuilder() *jwt.Builder {
	user := DefaultUser()
	return jwt.NewBuilder().
		Subject(user.ID.String()).
		Issuer(DefaultIssuerKeycloak).
		Audience([]string{DefaultAudience}).
		Claim("scope", DefaultScope)
}

func NewToken(t require.TestingT, options ...Option) jwt.Token {
	var o opts
	o.Apply(options...)

	b := NewTokenBuilder()

	if o.expired {
		t := time.Now().Add(-10 * time.Hour)
		t2 := t.Add(-5 * time.Hour)
		b = b.IssuedAt(t).
			Expiration(t2).
			NotBefore(t)
	} else {
		now := time.Now()
		b = b.IssuedAt(now).
			Expiration(time.Now().Add(10 * time.Hour)). //nolint:mnd
			NotBefore(now)
	}

	if o.mods != nil {
		o.mods(b)
	}

	tk, err := b.Build()
	require.NoError(t, err)

	return tk
}

// SignToken signs a JWT with `privateKey`.
func SignToken(t require.TestingT, tk jwt.Token, privateKey jwk.Key) string {
	alg, ok := privateKey.Algorithm()
	require.True(t, ok, "Algorithm is not set on private key.")

	jwtStr, err := jwt.Sign(tk, jwt.WithKey(alg, privateKey))
	require.NoError(t, err, "Could not sign token.")

	return string(jwtStr)
}

func (c *opts) Apply(options ...Option) {
	for _, f := range options {
		f(c)
	}
}

func WithModifications(f func(b *jwt.Builder)) Option {
	return func(o *opts) {
		if o.mods != nil {
			old := o.mods
			o.mods = func(b *jwt.Builder) {
				old(b)
				f(b)
			}
		} else {
			o.mods = f
		}
	}
}

func WithInvalid() Option {
	return func(o *opts) {
		o.expired = true
	}
}
