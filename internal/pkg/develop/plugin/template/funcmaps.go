package template

import "strings"

func FormatPackageName(name string) string {
	return strings.ReplaceAll(name, "-", "")
}
