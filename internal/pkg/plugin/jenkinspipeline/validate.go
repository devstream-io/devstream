package jenkinspipeline

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/types"
	"github.com/devstream-io/devstream/pkg/util/validator"
)

// setDefault config default fields for usage
func setDefault(options configmanager.RawOptions) (configmanager.RawOptions, error) {
	opts, err := newJobOptions(options)
	if err != nil {
		return nil, err
	}

	// config scm and projectRepo values
	projectRepo, err := opts.SCM.BuildRepoInfo()
	if err != nil {
		return nil, err
	}
	opts.ProjectRepo = projectRepo

	// set field value if empty
	if opts.Jenkins.Namespace == "" {
		opts.Jenkins.Namespace = "jenkins"
	}
	if opts.Pipeline.Job == "" {
		opts.Pipeline.Job = projectRepo.Repo
	}

	opts.CIConfig = opts.Pipeline.buildCIConfig(projectRepo, options.GetMapByKey("pipeline"))
	return types.EncodeStruct(opts)
}

func validate(options configmanager.RawOptions) (configmanager.RawOptions, error) {
	opts, err := newJobOptions(options)
	if err != nil {
		return nil, err
	}

	if err = validator.StructAllError(opts); err != nil {
		return nil, err
	}

	// check repo is valid
	if err := opts.ProjectRepo.CheckValid(); err != nil {
		log.Debugf("jenkins validate repo invalid: %+v", err)
		return nil, err
	}
	// check jenkins job name
	if err := opts.Pipeline.checkValid(); err != nil {
		log.Debugf("jenkins validate pipeline invalid: %+v", err)
		return nil, err
	}

	return options, nil
}
