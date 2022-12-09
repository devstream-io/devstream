package file

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

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

// ReadYamls reads file or files from dir whose suffix is yaml or yml
// and returns the content of the files without "---" separator
func ReadYamls(path string) ([]byte, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	var contents []byte
	if stat.IsDir() {
		filterYaml := func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			if filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml" {
				content, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				contents = append(contents, content...)
			}
			return nil
		}
		err = filepath.WalkDir(path, filterYaml)
	} else {
		contents, err = os.ReadFile(path)
	}

	if err != nil {
		return nil, err
	}

	return []byte(strings.ReplaceAll(string(contents), "\n---\n", "\n")), nil
}
