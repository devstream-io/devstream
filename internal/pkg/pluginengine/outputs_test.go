package pluginengine

import (
	"testing"
)

type stripOutputReferencePrefixAndSuffixTest struct {
	arg      string
	expected string
}

var stripOutputReferencePrefixAndSuffixTests = []stripOutputReferencePrefixAndSuffixTest{
	{"${{a.b.c.d}}", "a.b.c.d"},
	{"${{ a.b.c.d }}", "a.b.c.d"},
	{"${{  a.b.c.d  }}", "a.b.c.d"},
	{"${{  a.b.c  }}", "a.b.c"},
}

func TestStripOutputReferencePrefixAndSuffix(t *testing.T) {
	for _, test := range stripOutputReferencePrefixAndSuffixTests {
		if got := stripOutputReferencePrefixAndSuffix(test.arg); got != test.expected {
			t.Errorf("Output %s not equal to expected %s. Input: %s", got, test.expected, test.arg)
		}
	}
}

type isValidOutputsReferenceTest struct {
	arg      string
	expected bool
}

var isValidOutputsReferenceTests = []isValidOutputsReferenceTest{
	{"${{a.b.c.d}}", true},
	{"${{ a.b.c.d }}", true},
	{"${{  a.b.c.d  }}", true},
	{"${{ a.b.c }}", false},
	{"${{ a.b.c.d.e }}", true},
	{"${ a.b.c.d.e }", false},
	{"{{ a.b.c.d.e }}", false},
}

func TestIsValidOutputsReference(t *testing.T) {
	for _, test := range isValidOutputsReferenceTests {
		if got := isValidOutputsReferenceFormat(test.arg); got != test.expected {
			t.Errorf("Output %t not equal to expected %t. Input: %s", got, test.expected, test.arg)
		}
	}
}

type getToolNamePluginKindAndOutputReferenceKeyTest struct {
	arg       string
	expected1 string
	expected2 string
	expected3 string
}

var getToolNamePluginKindAndOutputReferenceKeyTests = []getToolNamePluginKindAndOutputReferenceKeyTest{
	{"name.kind.outputs.key", "name", "kind", "key"},
}

func TestGetToolNamePluginKindAndOutputReferenceKey(t *testing.T) {
	for _, test := range getToolNamePluginKindAndOutputReferenceKeyTests {
		if got1, got2, got3 := getToolNamePluginKindAndOutputReferenceKey(test.arg); got1 != test.expected1 || got2 != test.expected2 || got3 != test.expected3 {
			t.Errorf("Output %s, %s, %s not equal to expected %s, %s, %s.", got1, got2, got3, test.expected1, test.expected2, test.expected3)
		}
	}
}
