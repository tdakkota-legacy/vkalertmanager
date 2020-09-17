package emitter

import (
	"fmt"
	"text/template"

	"github.com/tdakkota/vkalertmanager/pkg/hook"
	"github.com/valyala/bytebufferpool"
)

const defaultTemplate = `
{{ range .Alerts }}
{{ if eq .Status "firing"}}🔥 <b>{{ .Status | toUpper }}</b> 🔥{{ else }}<b>{{ .Status | toUpper }}</b>{{ end }}
<b>{{ .Labels.alertname }}</b>
{{ if .Annotations.message }}
{{ .Annotations.message }}
{{ end }}
{{ if .Annotations.summary }}
{{ .Annotations.summary }}
{{ end }}
{{ if .Annotations.description }}
{{ .Annotations.description }}
{{ end }}
<b>Duration:</b> {{ duration .StartsAt .EndsAt }}{{ if ne .Status "firing"}}
<b>Ended:</b> {{ .EndsAt | since }}{{ end }}
{{ end }}
`

func Default() Template {
	t, err := Parse(defaultTemplate)
	if err != nil {
		panic(err)
	}

	return t
}

func ParseFiles(path ...string) (Template, error) {
	t, err := template.New("msg").ParseFiles(path...)
	if err != nil {
		return Template{}, fmt.Errorf("failed to parse template: %w", err)
	}

	return NewTemplate(t), nil
}

func Parse(s string) (Template, error) {
	t, err := template.New("msg").Parse(s)
	if err != nil {
		return Template{}, fmt.Errorf("failed to parse template: %w", err)
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
	if err != nil {
		return "", err
	}
	return b.String(), nil
}
