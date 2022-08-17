package template

// the interfaces and functions in the file is used to support better chain calls

type (
	WithGetter interface {
		String() (string, error)
		AddProcessor(processor Processor) WithGetter
		SetRender(render RenderFunc) WithRender
	}

	WithRender interface {
		Render() (string, error)
	}
)

func (r *render) SetContentGetter(getter ContentGetter) WithGetter {
	return r.setContentGetter(getter)
}

func (r *rendererWithGetter) AddProcessor(processor Processor) WithGetter {
	return r.addProcessor(processor)
}

func (r *rendererWithGetter) SetRender(render RenderFunc) WithRender {
	return r.setRender(render)
}

func (r *rendererWithRender) Render() (string, error) {
	return r.doRender()
}
