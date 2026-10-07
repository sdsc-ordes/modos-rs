package middlewear

import (
	"context"
	"strings"

	"github.com/sdsc-ordes/modos-rs/components/kms/internal/config"
	mdJwt "github.com/sdsc-ordes/modos-rs/components/kms/internal/jwt"
	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/auth"
	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/errors"
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
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
	) (any, error) {
		token, err := extractToken(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "missing token")
		}

		cl := mdJwt.NewClaims(&oidcCfg.ClaimBucketPermissions)
		// FIXME:
		// If we know which RPC is called when entering this middlewear:
		// We need to pass here a hardcoded map of `RPC function name -> scopes needed`
		// or do it in the handler functions them self by inspecting `cl`.
		var allowedScopes []string = nil
		err = auth.ValidateJWT(ctx, validator, token, allowedScopes, cl)
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
