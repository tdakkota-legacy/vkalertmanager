package template

import "sync"

const defaultTemplateText = `
{{ range .Alerts }}
{{ if eq .Status "firing"}}🔥 {{ .Status | toUpper }} 🔥{{ else }}{{ .Status | toUpper }}{{ end }}
{{ .Labels.alertname }}
{{ if .Annotations.message }}
{{ .Annotations.message }}
{{ end }}
{{ if .Annotations.summary }}
{{ .Annotations.summary }}
{{ end }}
{{ if .Annotations.description }}
{{ .Annotations.description }}
{{ end }}
Duration: {{ duration .StartsAt .EndsAt }}{{ if ne .Status "firing"}}
Ended: {{ .EndsAt | since }}{{ end }}
{{ end }}
`

var templateOnce = &sync.Once{}
var defaultTemplate *Template

func Default() *Template {
	templateOnce.Do(func() {
		t, err := Parse(defaultTemplateText)
		if err != nil {
			panic(err)
		}
		defaultTemplate = t
	})

	return defaultTemplate
}
