package graphql_test

import (
	"testing"

	gqltiming "github.com/CodeWarrior-debug/perspectize/backend/pkg/graphql"
	"github.com/stretchr/testify/assert"
)

func TestOperationTimer_ReturnsMiddleware(t *testing.T) {
	mw := gqltiming.OperationTimer()
	assert.NotNil(t, mw)
}
