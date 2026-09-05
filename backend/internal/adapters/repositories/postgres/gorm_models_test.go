package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGormModelTableNames(t *testing.T) {
	assert.Equal(t, "users", UserModel{}.TableName())
	assert.Equal(t, "categories", CategoryModel{}.TableName())
	assert.Equal(t, "content", ContentModel{}.TableName())
	assert.Equal(t, "perspectives", PerspectiveModel{}.TableName())
}
