package tree

import "strings"

// testdir/
// ├── dir1/
// │   ├── file1.go
// │   └── dir2/
// │       ├── file1.go
// │       └── file2.go
// └── FILE3.md
func ParseTree(treeText string) *TreeNode {
	lines := strings.Split(treeText, "\n")
	rootLine := strings.TrimSpace(lines[0])
	rootName := strings.TrimSuffix(rootLine, "/")
	root := NewTreeNode(rootName, true)
	stack := []*TreeNode{root}

	for _, line := range lines[1:] {
		// dir1 as an example, the indent is 3+2=5 in "├── dir1/".
		indent := strings.LastIndex(line, "──") + 2
		// line[indent:] apply to dir1/, the result is " dir1/"
		// name is "dir1/"
		name := strings.TrimSpace(line[indent:])
		isDir := strings.HasSuffix(name, "/")
		if isDir {
			name = strings.TrimSuffix(name, "/")
		}

		node := NewTreeNode(name, isDir)
		parent := stack[indent/4]
		parent.AddChild(node)

		if isDir {
			stack = append(stack, node)
		} else {
			stack = append(stack[:indent/4+1], node)
		}
	}
	return root
}
