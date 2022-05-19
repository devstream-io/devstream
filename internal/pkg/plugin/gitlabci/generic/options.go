package generic

// Options is the struct for configurations of the gitlabci-generic plugin.
type Options struct {
	PathWithNamespace string `validate:"required"`
	Branch            string `validate:"required"`
	TemplateURL       string `validate:"required"`
	TemplateVariables map[string]interface{}
}
