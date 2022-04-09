package pluginengine

import (
	"testing"
)

type getToolNamePluginOutputKeyTest struct {
	arg    string
	match  bool
	name   string
	plugin string
	key    string
}

var getToolNamePluginOutputKeyTests = []getToolNamePluginOutputKeyTest{
	{"${{a.b.outputs.d}}", true, "a", "b", "d"},
	{"prefix${{a.b.outputs.d}}suffix", true, "a", "b", "d"},
	{"${{  a.b.outputs.d  }}", true, "a", "b", "d"},
	{"${{  a.b.c  }}", false, "", "", ""},
}

func TestGetToolNamePluginOutputKey(t *testing.T) {
	for _, test := range getToolNamePluginOutputKeyTests {
		match, name, plugin, key := getToolNamePluginOutputKey(test.arg)
		if match != test.match || name != test.name || plugin != test.plugin || key != test.key {
			t.Errorf(
				"Output %t, %s, %s, %s not equal to expected %t, %s, %s, %s. Input: %s.",
				match, name, plugin, key,
				test.match, test.name, test.plugin, test.key,
				test.arg,
			)
		}
	}
}

type replaceOutputKeyWithValueTest struct {
	arg1     string
	arg2     string
	expected string
}

var replaceOutputKeyWithValueTests = []replaceOutputKeyWithValueTest{
	{"${{a.b.outputs.d}}", "value", "value"},
	{"prefix/${{ a.b.outputs.d }}/suffix", "value", "prefix/value/suffix"},
	{"${{a.b.c}}", "value", "${{a.b.c}}"},
}

func TestIsValidOutputsReference(t *testing.T) {
	for _, test := range replaceOutputKeyWithValueTests {
		if got := replaceOutputKeyWithValue(test.arg1, test.arg2); got != test.expected {
			t.Errorf("Output %s not equal to expected %s. Input: %s, %s.", got, test.expected, test.arg1, test.arg2)
		}
	}
}
