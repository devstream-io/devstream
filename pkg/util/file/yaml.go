package file

import (
	"bytes"
	"fmt"

	yamlUtil "github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"

	"github.com/devstream-io/devstream/pkg/util/pkgerror"
)

// YamlSequenceNode is yaml sequenceNode
type YamlSequenceNode struct {
	StrOrigin string
	StrArray  []string
}

func MergeYamlNode(dst *YamlSequenceNode, src *YamlSequenceNode) *YamlSequenceNode {
	if dst == nil {
		return src
	} else if src == nil {
		return dst
	}
	dst.StrOrigin = fmt.Sprintf("%s\n%s", dst.StrOrigin, src.StrOrigin)
	dst.StrArray = append(dst.StrArray, src.StrArray...)
	return dst
}

// IsEmpty check node fields are empty
func (n *YamlSequenceNode) IsEmpty() bool {
	return n.StrOrigin == "" && len(n.StrArray) == 0
}

// GetYamlNodeArrayByPath get element from yaml content by yaml path
// return string array for elements
func GetYamlNodeArrayByPath(content []byte, path string) (*YamlSequenceNode, error) {
	// 1. get ast node from yaml content
	node, err := getYamlAstNode(content, path)
	if err != nil {
		if pkgerror.CheckErrorMatchByMessage(err, "node not found") {
			return nil, nil
		}
		return nil, fmt.Errorf("yaml parse path[%s] failed:%w", path, err)
	}
	// 2. transfer node to sequence node
	seqNode, ok := node.(*ast.SequenceNode)
	if !ok {
		return nil, fmt.Errorf("yaml parse path[%s] is not valid sequenceNode", string(content))
	}
	var nodeArray = make([]string, 0)
	for _, sn := range seqNode.Values {
		nodeArray = append(nodeArray, sn.String())
	}
	y := &YamlSequenceNode{
		StrOrigin: node.String(),
		StrArray:  nodeArray,
	}
	return y, nil
}

// GetYamlNodeStrByPath get element from yaml content by yaml path
// return string format of node
func GetYamlNodeStrByPath(content []byte, path string) (string, error) {
	// 1. get ast node from yaml content
	node, err := getYamlAstNode(content, path)
	if err != nil {
		if pkgerror.CheckErrorMatchByMessage(err, "node not found") {
			return "", nil
		}
		return "", err
	}
	return node.String(), nil
}

func getYamlAstNode(content []byte, path string) (ast.Node, error) {
	appsYamlPath, err := yamlUtil.PathString(path)
	if err != nil {
		return nil, fmt.Errorf("yaml generate path failed: %w", err)
	}
	node, err := appsYamlPath.ReadNode(bytes.NewBuffer(content))
	if err != nil {
		return nil, fmt.Errorf("yaml read node failed: %w", err)
	}
	return node, nil
}
