package emitter

import (
	"text/template"

	"github.com/tdakkota/vkalertmanager/pkg/hook"
	"github.com/valyala/bytebufferpool"
)

func Parse(s string) (Template, error) {
	t, err := template.New("msg").Parse(s)
	if err != nil {
		return Template{}, err
	}

	return NewTemplate(t), nil
}

type Template struct {
	template *template.Template
	pool     *bytebufferpool.Pool
}

func NewTemplate(template *template.Template) Template {
	return Template{
		template: template,
		pool:     &bytebufferpool.Pool{},
	}
}

func (t Template) Produce(m hook.Message) (string, error) {
	b := t.pool.Get()
	defer t.pool.Put(b)

	err := t.template.Execute(b, m)
	return b.String(), err
}
