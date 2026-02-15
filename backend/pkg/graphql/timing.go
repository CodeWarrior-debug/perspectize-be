package graphql

import (
	"context"
	"log/slog"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

// OperationTimer returns a gqlgen OperationMiddleware that logs
// the operation name and duration for every GraphQL request.
func OperationTimer() graphql.OperationMiddleware {
	return func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		start := time.Now()
		oc := graphql.GetOperationContext(ctx)

		rh := next(ctx)

		operationName := "anonymous"
		if oc != nil && oc.OperationName != "" {
			operationName = oc.OperationName
		}

		slog.InfoContext(ctx, "graphql",
			"operation", operationName,
			"duration_ms", time.Since(start).Milliseconds(),
		)

		return rh
	}
}
