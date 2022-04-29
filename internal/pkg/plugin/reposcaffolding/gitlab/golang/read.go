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
	c, err := gitlab.NewClient()
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
	res["owner"] = project.Owner.Username
	res["org"] = project.Owner.Organization
	res["repoName"] = project.Name

	outputs := make(map[string]interface{})
	if opts.Owner == "" {
		outputs["owner"] = opts.Owner
	} else {
		outputs["owner"] = project.Owner.Username
	}
	if opts.Org == "" {
		outputs["org"] = opts.Org
	} else {
		outputs["org"] = project.Owner.Organization
	}
	outputs["repo"] = project.Name
	outputs["repoURL"] = project.HTTPURLToRepo
	res["outputs"] = outputs

	return res, nil
}
