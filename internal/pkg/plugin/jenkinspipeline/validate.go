package jenkinspipeline

import (
	_ "embed"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/cifile"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/types"
	"github.com/devstream-io/devstream/pkg/util/validator"
)

//go:embed tpl/Jenkinsfile_offline.tpl
var offlineJenkinsScript string

// setJenkinsDefault config default fields for usage
func setJenkinsDefault(options configmanager.RawOptions) (configmanager.RawOptions, error) {
	opts, err := newJobOptions(options)
	if err != nil {
		return nil, err
	}
	// set project and ci default
	projectRepo, err := opts.SCM.BuildRepoInfo()
	if err != nil {
		return nil, err
	}
	opts.ProjectRepo = projectRepo
	// if jenkins is offline, just use offline Jenkinsfile
	if opts.needOfflineConfig() {
		opts.CIFileConfig = &cifile.CIFileConfig{
			Type: ciType,
			ConfigContentMap: map[string]string{
				"Jenkinsfile": offlineJenkinsScript,
			},
			Vars: opts.Pipeline.GenerateCIFileVars(projectRepo),
		}
	} else {
		opts.CIFileConfig = opts.Pipeline.BuildCIFileConfig(ciType, projectRepo)
	}
	// set field value if empty
	if opts.Jenkins.Namespace == "" {
		opts.Jenkins.Namespace = "jenkins"
	}
	if opts.JobName == "" && opts.ProjectRepo != nil {
		opts.JobName = jenkinsJobName(opts.ProjectRepo.Repo)
	}
	return types.EncodeStruct(opts)
}

// validateJenkins will validate jenkins jobName field and jenkins field
func validateJenkins(options configmanager.RawOptions) (configmanager.RawOptions, error) {
	opts, err := newJobOptions(options)
	if err != nil {
		return nil, err
	}

	if err = validator.StructAllError(opts); err != nil {
		return nil, err
	}

	if err := opts.ProjectRepo.CheckValid(); err != nil {
		log.Debugf("github action validate repo invalid: %+v", err)
		return nil, err
	}

	// check jenkins job name
	if err := opts.JobName.checkValid(); err != nil {
		log.Debugf("jenkins validate pipeline invalid: %+v", err)
		return nil, err
	}
	return options, nil
}
