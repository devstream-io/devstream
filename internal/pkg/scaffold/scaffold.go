package scaffold

func Scaffold(tree string) error {
	return nil
}

// parseTreeToList parses the directory tree represented by a string and
// returns a list of directories and files.
// for example input:"
// testdir/
// ├── dir1/
// │   ├── file1.go
// │   └── dir2/
// │       ├── file1.go
// │       └── file2.go
// └── FILE3.md
// "
// output:
//
//	[]string{
//				"testdir/",
//				"testdir/dir1/",
//				"testdir/dir1/file1.go",
//				"testdir/dir1/dir2/",
//				"testdir/dir1/dir2/file1.go",
//				"testdir/dir1/dir2/file2.go",
//	}
//func parseTreeToList(tree string) ([]string, error) {
//	tree = strings.ReplaceAll(tree, "\r\n", "\n")
//	tree = strings.ReplaceAll(tree, "\r", "\n")
//	tree = strings.ReplaceAll(tree, "\t", "")
//	tree = strings.ReplaceAll(tree, "\n\n", "\n")
//	tree = strings.ReplaceAll(tree, "─", "-")
//	tree = strings.ReplaceAll(tree, "│", "|")
//
//	lines := strings.Split(tree, "\n")
//	var result []string
//
//	for _, line := range lines {
//		line = strings.TrimSpace(line)
//
//		if len(line) == 0 {
//			continue
//		}
//
//		level := 0
//		for ; strings.HasPrefix(line, "  "); level++ {
//			line = line[2:]
//		}
//
//		prefix := strings.Repeat("  ", level)
//
//		if strings.HasSuffix(line, "/") {
//			result = append(result, prefix+line)
//		} else {
//			result = append(result, prefix+line)
//		}
//
//	}
//
//	return result, nil
//}
//
//func ParseTree(tree string) ([]string, error) {
//	var result []string
//
//	lines := strings.Split(tree, "\n")
//	for _, line := range lines {
//		if strings.TrimSpace(line) == "" {
//			continue
//		}
//
//		parts := strings.Split(line, " ")
//		level := strings.Count(parts[0], "│") + strings.Count(parts[0], "└") + strings.Count(parts[0], "├")
//		name := strings.TrimSpace(parts[len(parts)-1])
//
//		// Construct the full path by combining the names of all parent directories
//		path := ""
//		for i := 1; i < level; i++ {
//			if len(result) < i {
//				return nil, fmt.Errorf("parent directory not found: %s", line)
//			}
//			path += result[i-1] + "/"
//		}
//		path += name
//
//		result = append(result, path)
//	}
//
//	return result, nil
//}
