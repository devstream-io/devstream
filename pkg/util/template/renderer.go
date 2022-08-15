package template

type (
	// ContentGetter gets content from any source
	ContentGetter interface {
		GetContent() ([]byte, error)
	}

	// Processor process content before render
	Processor interface {
		Process([]byte) ([]byte, error)
	}

	// RenderInf render content to string
	RenderInf interface {
		Render([]byte) (string, error)
	}
)

type (

	// the following three structs implement the "State Pattern"
	// e.g. if render calls any function which set the getter, it turns into rendererWithGetter

	// why I use this pattern?
	// 1. because it is easy to prevent the caller from forget to set the getter.
	// 2. and it is friendly to the code hints.
	//	e.g. if you calls template.New() and get a render struct,
	//       IDE will only show you the methods which set the getter.
	render struct {
	}

	rendererWithGetter struct {
		getter     ContentGetter
		processors []Processor
	}

	rendererWithRender struct {
		getter     ContentGetter // mandatory
		processors []Processor   // optional
		render     RenderInf     // mandatory
	}
)

func New() *render {
	return &render{}
}

func (r *render) SetContentGetter(getter ContentGetter) *rendererWithGetter {
	return &rendererWithGetter{
		getter: getter,
	}
}

func (r *rendererWithGetter) AddProcessor(processor Processor) *rendererWithGetter {
	return &rendererWithGetter{
		getter:     r.getter,
		processors: append(r.processors, processor),
	}
}

func (r *rendererWithGetter) SetRender(render RenderInf) *rendererWithRender {
	return &rendererWithRender{
		getter:     r.getter,
		processors: r.processors,
		render:     render,
	}
}

// Render gets the content, process the content, render and returns the result string
func (c *rendererWithRender) Render() (string, error) {
	// 1. get content
	content, err := c.getter.GetContent()
	if err != nil {
		return "", err
	}

	// 2. process content
	for _, processor := range c.processors {
		content, err = processor.Process(content)
		if err != nil {
			return "", err
		}
	}

	// 3. render content
	return c.render.Render(content)
}

// String returns the string directly, without rendering
func (c *rendererWithGetter) String() (string, error) {
	// 1. get content
	content, err := c.getter.GetContent()
	if err != nil {
		return "", err
	}

	// 2. process content
	for _, processor := range c.processors {
		content, err = processor.Process(content)
		if err != nil {
			return "", err
		}
	}

	return string(content), nil
}
