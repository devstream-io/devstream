package golang

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	rs "github.com/devstream-io/devstream/internal/pkg/plugin/common/reposcaffolding"
	"github.com/devstream-io/devstream/pkg/util/gitlab"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Read(options map[string]interface{}) (map[string]interface{}, error) {
	var opts rs.Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}

	if errs := rs.Validate(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return nil, fmt.Errorf("opts are illegal")
	}

	return buildReadState(&opts)
}

func buildReadState(opts *rs.Options) (map[string]interface{}, error) {
	c, err := gitlab.NewClient(gitlab.WithBaseURL(opts.BaseURL))
	if err != nil {
		return nil, err
	}

	project, err := c.DescribeProject(opts.PathWithNamespace)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, nil
	}

	res := make(map[string]interface{})
	outputs := make(map[string]interface{})

	log.Debugf("GitLab Project is: %#v\n", project)

	if project.Owner != nil {
		log.Debugf("GitLab project owner is: %#v.\n", project.Owner)
		res["owner"] = project.Owner.Username
		res["org"] = project.Owner.Organization
		outputs["owner"] = project.Owner.Username
		outputs["org"] = project.Owner.Organization
	} else {
		res["owner"] = opts.Owner
		res["org"] = opts.Org
		outputs["owner"] = opts.Owner
		outputs["org"] = opts.Org
	}
	res["repoName"] = project.Name
	outputs["repo"] = project.Name
	outputs["repoURL"] = project.HTTPURLToRepo
	res["outputs"] = outputs

	return res, nil
}
