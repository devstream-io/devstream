//go:build ignore
// +build ignore

// This program is run via "go generate" to generate the code.

package main

import (
	_ "embed"
	"flag"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/devstream-io/devstream/pkg/util/file"
	templateUtil "github.com/devstream-io/devstream/pkg/util/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Option struct {
	Plugins []string
	Dir     string
	Path    string
	Package string
	Funcs   template.FuncMap
}

const (
	templatesSrc = "../../../../examples"
	templatesDst = "templates"
)

func main() {
	srcDir := flag.String("dir", "plugins", "source directory for yaml files")
	packageName := flag.String("pkg", "config", "package name")
	flag.Parse()

	if err := copyTemplates(templatesSrc, templatesDst); err != nil {
		log.Fatal(err)
	}

	plugins, err := getYamlFiles(*srcDir)
	if err != nil {
		log.Fatal(err)
	}

	generate(&Option{
		Plugins: plugins,
		Dir:     *srcDir,
		Path:    "embed_gen.go",
		Package: *packageName,
		Funcs: template.FuncMap{
			"UpperCamelCase": UpperCamelCase,
		},
	})
}

//go:embed gen_embed.tpl
var templateCode string

func copyTemplates(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".yaml") {
			return nil
		}
		return file.CopyFile(path, filepath.Join(dst, filepath.Base(path)))
	})
}

// generate generates the code for Option `o` into a file named by `o.Path`.
func generate(o *Option) {
	content, err := templateUtil.NewRenderClient(&templateUtil.TemplateOption{
		Name:    "gen_embed",
		FuncMap: o.Funcs,
	}, templateUtil.ContentGetter).Render(templateCode, o)

	if err != nil {
		log.Fatal("template Execute error:", err)
	}

	formatted, err := format.Source([]byte(content))
	if err != nil {
		log.Fatal("format:", err)
		log.Fatal("template Execute error:", err)
	}

	if err := os.WriteFile(o.Path, formatted, 0644); err != nil {
		log.Fatal("writeFile:", err)
	}

}

// getYamlFiles returns a list of YAML files' names in the given directory.
func getYamlFiles(dir string) ([]string, error) {

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var fileNames []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if !strings.HasSuffix(file.Name(), ".yaml") {
			continue
		}
		fileNames = append(fileNames, strings.TrimSuffix(file.Name(), ".yaml"))
	}

	return fileNames, nil
}

// UpperCamelCase returns a string with the first letter in upper case.
func UpperCamelCase(s string) string {
	s = strings.Replace(s, "-", " ", -1)
	s = cases.Title(language.English).String(s)
	return strings.Replace(s, " ", "", -1)
}
