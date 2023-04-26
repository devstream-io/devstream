package tree

import "fmt"

type TreeNode struct {
	Name     string
	IsDir    bool
	Children []*TreeNode
}

func NewTreeNode(name string, isDir bool) *TreeNode {
	return &TreeNode{
		Name:     name,
		IsDir:    isDir,
		Children: []*TreeNode{},
	}
}

func (t *TreeNode) AddChild(child *TreeNode) {
	t.Children = append(t.Children, child)
}

func (t *TreeNode) PrintTree(prefix string) {
	if t.IsDir {
		fmt.Println(prefix + t.Name + "/")
	} else {
		fmt.Println(prefix + t.Name)
	}
	for _, child := range t.Children {
		child.PrintTree(prefix + "    ")
	}
}
