package interceptor

import (
	"context"
	"strings"

	"github.com/abhishek622/moviedock/pkg/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey string

const UserIDKey contextKey = "userID"

func UnaryAuthInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
	}

	authHeaders := md["authorization"]
	if len(authHeaders) == 0 {
		return nil, status.Error(codes.Unauthenticated, "authorization token is not provided")
	}

	authHeader := authHeaders[0]
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, status.Error(codes.Unauthenticated, "invalid authorization header format")
	}

	tokenString := parts[1]
	claims, err := auth.ValidateToken(tokenString)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid or expired token")
	}

	ctxWithUser := context.WithValue(ctx, UserIDKey, claims.UserID)

	return handler(ctxWithUser, req)
}
