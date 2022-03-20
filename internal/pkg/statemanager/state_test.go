package statemanager

import "testing"

type generateStateKeyByToolNameAndPluginKindTest struct {
	arg1     string
	arg2     string
	expected StateKey
}

var generateStateKeyByToolNameAndPluginKindTests = []generateStateKeyByToolNameAndPluginKindTest{
	{"name", "kind", "name_kind"},
}

func TestGenerateStateKeyByToolNameAndPluginKind(t *testing.T) {
	for _, test := range generateStateKeyByToolNameAndPluginKindTests {
		if got := GenerateStateKeyByToolNameAndPluginKind(test.arg1, test.arg2); got != test.expected {
			t.Errorf("Output %s not equal to expected %s. Input: %s, %s.", got, test.expected, test.arg1, test.arg2)
		}
	}
}
