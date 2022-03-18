package plugin

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"text/template"

	pluginTpl "github.com/merico-dev/stream/internal/pkg/develop/plugin/template"
	"github.com/merico-dev/stream/pkg/util/log"
)

type Plugin struct {
	Name string
}

func NewPlugin(name string) *Plugin {
	return &Plugin{
		Name: name,
	}
}

// RenderTplFiles takes specified data that the templates needed,
// then render TplFiles to "Files" and return it as []File.
func (p *Plugin) RenderTplFiles() ([]pluginTpl.File, error) {
	retFiles := make([]pluginTpl.File, 0)
	count := len(pluginTpl.TplFiles)
	log.Debugf("Template files count: %d.", count)

	for i, tplFile := range pluginTpl.TplFiles {
		log.Debugf("Render process: %d/%d.", i+1, count)
		file, err := p.renderTplFile(&tplFile)
		if err != nil {
			return nil, err
		}
		log.Debugf("File: %v.", *file)
		retFiles = append(retFiles, *file)
	}

	return retFiles, nil
}

// RenderTplFile takes specified data that the template needed,
// then render one TplFile to File and return it.
func (p *Plugin) renderTplFile(tplFile *pluginTpl.TplFile) (*pluginTpl.File, error) {
	name, err := p.renderTplString(tplFile.NameTpl)
	if err != nil {
		return nil, err
	}
	dir, err := p.renderTplString(tplFile.DirTpl)
	if err != nil {
		return nil, err
	}
	content, err := p.renderTplString(tplFile.ContentTpl)
	if err != nil {
		return nil, err
	}

	return &pluginTpl.File{
		Name:    name,
		Dir:     dir,
		Content: content,
	}, nil
}

// renderTplString get the template string and the data object,
// then render it and return the rendered string.
func (p *Plugin) renderTplString(tplStr string) (string, error) {
	if tplStr == "" {
		return "", nil
	}

	t, err := template.New("default").Parse(tplStr)
	if err != nil {
		log.Debugf("Template parse failed: %s.", err)
		log.Debugf("Template content: %s.", tplStr)
		return "", err
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, *p)
	if err != nil {
		log.Debugf("Template execute failed: %s.", err)
		log.Debugf("Template content: %s.", tplStr)
		log.Debugf("Data object: %v.", *p)
		return "", err
	}

	return buf.String(), nil
}

// PersistFiles gets the []pluginTpl.File, for each File:
// call the persistFile() method to deal with.
func (p *Plugin) PersistFiles(files []pluginTpl.File) error {
	fileCount := len(files)
	log.Debugf("There are %d files wait to persist.", fileCount)
	for i, file := range files {
		log.Debugf("Persist process: %d/%d", i+1, fileCount)
		if err := p.persistFile(&file); err != nil {
			log.Errorf("Failed to persist: %s/%s", file.Dir, file.Name)
			return err
		}
	}

	return nil
}

// persistFile gets the *pluginTpl.File, then do the following:
// 1. mkdir the File.Dir
// 2. create the File.Name file
// 3. write the File.Content into the File.Name file
func (p *Plugin) persistFile(file *pluginTpl.File) error {
	// mkdir the File.Dir
	if err := os.MkdirAll(file.Dir, 0755); err != nil {
		log.Debugf("Failed to create directory: %s.", file.Dir)
	}
	log.Debugf("Directory created: %s.", file.Dir)

	// create the File.Name file and write the File.Content into the File.Name file
	filePath := path.Join(file.Dir, file.Name)
	if err := os.WriteFile(filePath, []byte(file.Content), 0644); err != nil {
		log.Debugf("Failed to write content to the file: %s.", err)
	}
	log.Debugf("File %s has been created.", filePath)

	return nil
}

func (p *Plugin) PrintHelpInfo() {
	help := `
The DevStream PMC (project management committee) sincerely thank you for your devotion and enthusiasm in creating new plugins!

To make the process easy as a breeze, DevStream(dtm) has generated some templated source code files for you to flatten the learning curve and reduce manual copy-paste.
In the generated templates, dtm has left some special marks in the format of "TODO(dtm)".
Please look for these TODOs by global search. Once you find them, you will know what to do with them. Also, please remember to check our documentation on creating a new plugin:

**README_when_create_plugin.md**

Source code files created.

Happy hacking, buddy!
Please give us feedback through GitHub issues if you encounter any difficulties. We guarantee that you will receive unrivaled help from our passionate community!
`
	fmt.Println(help)
}
