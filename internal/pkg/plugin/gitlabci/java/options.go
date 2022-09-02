package java

import (
	"errors"
	"fmt"
	"html/template"
	"strings"

	"github.com/spf13/viper"

	"github.com/devstream-io/devstream/pkg/util/scm/git"
	"github.com/devstream-io/devstream/pkg/util/scm/gitlab"
)

type RegistryType string

const (
	RegistryDockerhub RegistryType = "dockerhub"
	RegistryHarbor    RegistryType = "harbor"
)

// Options is the struct for configurations of the gitlabci-java plugin.
type Options struct {
	PathWithNamespace string   `validate:"required"`
	Branch            string   `validate:"required"`
	BaseURL           string   `validate:"omitempty,url"`
	Package           *Package `validate:"required"`
	Build             *Build   `validate:"required"`
	Deploy            *Deploy  `validate:"required"`
}

type BaseOption struct {
	Enable        bool
	Image         string
	Tags          string
	AllowedBranch []string
}

type Package struct {
	*BaseOption   `validate:"required"`
	ScriptCommand []string
}

type Build struct {
	*BaseOption   `validate:"required"`
	Registry      *Registry
	ImageName     string
	ScriptCommand []string
}

type Registry struct {
	Type     RegistryType
	Username string
}

type Deploy struct {
	*BaseOption   `validate:"required"`
	ScriptCommand []template.HTML
	K8sAgentName  string
}

// Set options with default value
func (o *Options) complete() error {
	if o.Package.Enable {
		o.Package.setup()
	}

	if o.Build.Enable {
		if err := o.Build.setup(); err != nil {
			return err
		}
	}

	if o.Deploy.Enable {
		o.Deploy.setup(o)
	}

	return nil
}

func (p *Package) setup() {
	if p.Image == "" {
		p.Image = defaultMVNPackageJobImg
	}

	if len(p.ScriptCommand) == 0 {
		p.ScriptCommand = append(p.ScriptCommand, defaultMVNPackageJobScript)
	}

	if p.Tags == "" {
		p.Tags = defaultTags
	}

	if len(p.AllowedBranch) == 0 {
		p.AllowedBranch = append(p.AllowedBranch, "main")
	}
}

func (b *Build) setup() error {
	if b.Image == "" {
		b.Image = defaultDockerBuildJobImg
	}

	if len(b.ScriptCommand) == 0 {
		dockerhubToken := viper.GetString("dockerhub_token")
		if dockerhubToken == "" {
			return fmt.Errorf("DockerHub Token is empty")
		}

		defaultDuildScripts := []string{
			fmt.Sprintf("docker login -u %s -p %s", b.Registry.Username, dockerhubToken),
			fmt.Sprintf("docker build -t %s/%s:$CI_PIPELINE_ID .", b.Registry.Username, b.ImageName),
			fmt.Sprintf("docker push %s/%s:$CI_PIPELINE_ID", b.Registry.Username, b.ImageName),
		}

		b.ScriptCommand = append(b.ScriptCommand, defaultDuildScripts...)
	}

	if b.Tags == "" {
		b.Tags = defaultTags
	}

	if len(b.AllowedBranch) == 0 {
		b.AllowedBranch = append(b.AllowedBranch, "main")
	}

	return nil
}

func (d *Deploy) setup(opts *Options) {
	if d.Image == "" {
		d.Image = defaultK8sDeployJobImg
	}

	if len(d.ScriptCommand) == 0 {
		defalutDeployScripts := []template.HTML{
			"kubectl config get-contexts",
			template.HTML(fmt.Sprintf("- kubectl config use-context %s:%s", opts.PathWithNamespace, d.K8sAgentName)),
			"cd manifests",
			template.HTML(`sed -i "s/IMAGE_TAG/$CI_PIPELINE_ID/g" deployment.yaml`),
			"cat deployment.yaml",
			"kubectl apply -f deployment.yaml",
		}
		d.ScriptCommand = append(d.ScriptCommand, defalutDeployScripts...)
	}

	if d.Tags == "" {
		d.Tags = defaultTags
	}

	if len(d.AllowedBranch) == 0 {
		d.AllowedBranch = append(d.AllowedBranch, "main")
	}

}

func (opts *Options) newGitlabClient() (*gitlab.Client, error) {
	pathSplit := strings.Split(opts.PathWithNamespace, "/")
	if len(pathSplit) != 2 {
		return nil, errors.New("gitlabci generic not valid PathWithNamespace params")
	}
	repoInfo := &git.RepoInfo{
		Owner:   pathSplit[0],
		Repo:    pathSplit[1],
		Branch:  opts.Branch,
		BaseURL: opts.BaseURL,
	}
	return gitlab.NewClient(repoInfo)
}
