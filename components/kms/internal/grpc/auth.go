package grpc

import (
	"context"
	"log"
	"strings"

	"github.com/sdsc-ordes/modos-rs/components/kms/internal/config"
	mdJwt "github.com/sdsc-ordes/modos-rs/components/kms/internal/jwt"
	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/auth"
	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/errors"
	clog "gitlab.com/data-custodian/custodian/components/lib-common/pkg/log/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const authHeader = "authorization"
const KeyClaims ctxKey = 10000

// To make the key not collide with other packages,
// since the value and type must be matched.
type ctxKey int

func AuthenticationInterceptor(
	validator *auth.JWTVerifier,
	oidcCfg *config.OIDC,
) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
	) (any, error) {
		token, err := extractToken(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "missing token")
		}

		expectedScopes, err := ExpectedScopes(info.FullMethod)
		if err != nil {
			clog.ErrorEf(ctx, err, "Failed to determine allowed scopes for '%v'", info.FullMethod)
		}

		cl := mdJwt.NewClaims(&oidcCfg.ClaimBucketPermissions)
		err = auth.ValidateJWT(ctx, validator, token, expectedScopes, cl)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		return handler(
			context.WithValue(ctx, KeyClaims, cl),
			req)
	}
}

// Claims returns the claims extracted by the [AuthenticationInterceptor].
func Claims(ctx context.Context) *mdJwt.Claims {
	if c, ok := ctx.Value(KeyClaims).(*mdJwt.Claims); ok {
		return c
	}

	log.Panic(
		"Could not extract claim key from context. " +
			"Did you set the claims on the context?")

	return nil
}

func extractToken(ctx context.Context) (string, error) {
	var token string
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("no metadata found on context")
	}

	authData := md.Get(authHeader)
	if len(authData) == 0 {
		return "", errors.New("no auth. header '%v' found in meta data", authHeader)
	} else {
		token = strings.TrimPrefix(authData[0], "Bearer ")
	}

	if token == "" {
		return "", errors.New("no token provided in '%s' in meta data", authHeader)
	}

	return token, nil
}
