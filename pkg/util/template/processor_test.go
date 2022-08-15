package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddDotForVariablesInConfigNormal(t *testing.T) {
	res := AddDotForVariablesInConfig("[[varNameA]]")
	assert.Equal(t, "[[ .varNameA]]", res, "Adding dot for variable names passed.")
}

func TestAddDotForVariablesInConfigWithSpaces(t *testing.T) {
	res := AddDotForVariablesInConfig("[[  varNameA]]")
	assert.Equal(t, "[[ .varNameA]]", res, "Adding dot for variable names passed.")
}

func TestAddDotForVariablesInConfigWithTrailingSpaces(t *testing.T) {
	res := AddDotForVariablesInConfig("[[ varNameA  ]]")
	assert.Equal(t, "[[ .varNameA  ]]", res, "Adding dot for variable names passed.")
}

func TestAddDotForVariablesInConfigMultipleVars(t *testing.T) {
	res := AddDotForVariablesInConfig("[[ varNameA ]]/[[ varNameB ]]/[[ varNameC ]]")
	assert.Equal(t, "[[ .varNameA ]]/[[ .varNameB ]]/[[ .varNameC ]]", res, "Adding dot for variable names passed.")
}
