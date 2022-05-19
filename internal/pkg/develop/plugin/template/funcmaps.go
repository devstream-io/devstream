package template

import "strings"

func FormatPackageName(name string) string {
	if specialPluginName, isOk := SpecialPluginNameMap[name]; isOk {
		return specialPluginName.PackageName
	} else {
		// Remove suffixes like "-integ"
		name = strings.TrimSuffix(name, "-integ")
		return strings.ReplaceAll(name, "-", "")
	}
}

func FormatPackageDirName(name string) string {
	if specialPluginName, isOk := SpecialPluginNameMap[name]; isOk {
		return specialPluginName.DirName
	} else {
		// Remove suffixes like "-integ"
		name = strings.TrimSuffix(name, "-integ")
		return strings.ReplaceAll(name, "-", "")
	}
}
