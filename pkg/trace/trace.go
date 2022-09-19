package trace

import (
	"context"
	"github.com/google/uuid"
)

const TIDKey = "trace_id"

// Ctx will generate a UUID for tracing and logging.
// Returns a context with a trace ID to be used.
func Ctx() context.Context {
	return context.WithValue(context.Background(), TIDKey, uuid.New().String())
}
