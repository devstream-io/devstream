package configloader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddDotForVariablesInConfigNormal(t *testing.T) {
	res := addDotForVariablesInConfig("[[varNameA]]")
	assert.Equal(t, "[[ .varNameA]]", res, "Adding dot for variable names passed.")
}

func TestAddDotForVariablesInConfigWithSpaces(t *testing.T) {
	res := addDotForVariablesInConfig("[[  varNameA]]")
	assert.Equal(t, "[[ .varNameA]]", res, "Adding dot for variable names passed.")
}

func TestAddDotForVariablesInConfigWithTrailingSpaces(t *testing.T) {
	res := addDotForVariablesInConfig("[[ varNameA  ]]")
	assert.Equal(t, "[[ .varNameA  ]]", res, "Adding dot for variable names passed.")
}

func TestAddDotForVariablesInConfigMultipleVars(t *testing.T) {
	res := addDotForVariablesInConfig("[[ varNameA ]]/[[ varNameB ]]/[[ varNameC ]]")
	assert.Equal(t, "[[ .varNameA ]]/[[ .varNameB ]]/[[ .varNameC ]]", res, "Adding dot for variable names passed.")
}

func TestRenderConfigWithVariables(t *testing.T) {
	variables := map[string]interface{}{
		"varNameA": "A",
		"varNameB": "B",
	}
	result, err := renderConfigWithVariables("[[ .varNameA ]]/[[ .varNameB]]", variables)
	assert.Equal(t, err, nil)
	assert.Equal(t, string(result), "A/B")
}
