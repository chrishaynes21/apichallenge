package trace

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCtx(t *testing.T) {
	// Table test not needed, this test gets all branches
	t.Run("should return new context with trace id", func(t *testing.T) {
		gotCtx := Ctx()
		assert.NotEmpty(t, gotCtx.Value(TIDKey), "context trace id should not be empty")
	})
}
