package jenkins

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/types"
	"github.com/devstream-io/devstream/pkg/util/validator"
)

// SetJobDefaultConfig config default fields for usage
func SetJobDefaultConfig(options configmanager.RawOptions) (configmanager.RawOptions, error) {
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

	// set ci field field
	ciConfig, err := opts.Pipeline.buildCIConfig(projectRepo)
	if err != nil {
		return nil, err
	}
	opts.CIConfig = ciConfig
	return types.EncodeStruct(opts)
}

func ValidateJobConfig(options configmanager.RawOptions) (configmanager.RawOptions, error) {
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
