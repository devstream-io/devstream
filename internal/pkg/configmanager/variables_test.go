package configmanager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderConfigWithVariables(t *testing.T) {
	variables := map[string]interface{}{
		"varNameA": "A",
		"varNameB": "B",
	}
	result, err := renderConfigWithVariables("[[ varNameA ]]/[[ varNameB]]", variables)
	assert.Equal(t, err, nil)
	assert.Equal(t, string(result), "A/B")
}
