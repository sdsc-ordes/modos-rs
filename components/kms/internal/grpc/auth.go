package grpc

import (
	"context"
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

		cl := mdJwt.NewClaims(&oidcCfg.ClaimBucketPermissions)

		expectedScopes, err := ExpectedScopes(info.FullMethod)
		if err != nil {
			clog.ErrorEf(ctx, err, "Failed to determine allowed scopes for '%v'", info.FullMethod)
		}

		err = auth.ValidateJWT(ctx, validator, token, expectedScopes, cl)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		return handler(ctx, req)
	}
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
