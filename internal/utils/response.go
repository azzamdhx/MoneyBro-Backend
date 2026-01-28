package utils

import (
	"context"
	"errors"

	"github.com/vektah/gqlparser/v2/gqlerror"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrNotFound     = errors.New("not found")
	ErrBadRequest   = errors.New("bad request")
	ErrEmailExists  = errors.New("email already exists")
)

func GraphQLError(ctx context.Context, message string, code string) *gqlerror.Error {
	return &gqlerror.Error{
		Message: message,
		Extensions: map[string]interface{}{
			"code": code,
		},
	}
}

func UnauthorizedError(ctx context.Context) *gqlerror.Error {
	return GraphQLError(ctx, "Unauthorized", "UNAUTHENTICATED")
}

func ForbiddenError(ctx context.Context) *gqlerror.Error {
	return GraphQLError(ctx, "Forbidden", "FORBIDDEN")
}

func NotFoundError(ctx context.Context, resource string) *gqlerror.Error {
	return GraphQLError(ctx, resource+" not found", "NOT_FOUND")
}

func ValidationError(ctx context.Context, message string) *gqlerror.Error {
	return GraphQLError(ctx, message, "VALIDATION_ERROR")
}
