package base

type BasePipelineConfig map[string]string

type StepOperation interface {
	PreConfig() (BasePipelineConfig, error)
	Render() BasePipelineConfig
}

type BasePipeline struct {
	ImageRepo      StepOperation
	Dingtalk       StepOperation
	StepOperation  StepOperation
	ConfigFilePath string
}
