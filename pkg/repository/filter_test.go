package repository

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProjectFilterBuilder(t *testing.T) {
	builder := NewProjectFilterPatternBuilder(
		"myproject",
	)

	pattern := builder.Build()

	assert.Equal(t, pattern, "^myproject.*v([0-9]{3})$")
}

func TestProjectAndApplicationFilterBuilder(t *testing.T) {
	builder := NewProjectAndApplicationFilterPatternBuilder(
		"myproject",
		"myapp",
	)

	pattern := builder.Build()

	assert.Equal(t, pattern, "^myproject[\\-\\_]myapp\\-v([0-9]{3})$")
}
