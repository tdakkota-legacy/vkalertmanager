package template

import (
	"fmt"
	html "html/template"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/hako/durafmt"

	"github.com/tdakkota/vkalertmanager/pkg/hook"
	"github.com/valyala/bytebufferpool"
)

var DefaultFuncs = template.FuncMap{
	"toUpper": strings.ToUpper,
	"toLower": strings.ToLower,
	"title":   strings.Title,
	// join is equal to strings.Join but inverts the argument order
	// for easier pipelining in templates.
	"join": func(sep string, s []string) string {
		return strings.Join(s, sep)
	},
	"match": regexp.MatchString,
	"safeHtml": func(text string) html.HTML {
		return html.HTML(text)
	},
	"reReplaceAll": func(pattern, repl, text string) string {
		re := regexp.MustCompile(pattern)
		return re.ReplaceAllString(text, repl)
	},
	"stringSlice": func(s ...string) []string {
		return s
	},
	"since": func(t time.Time) string {
		return durafmt.Parse(time.Since(t)).String()
	},
	"duration": func(start time.Time, end time.Time) string {
		return durafmt.Parse(end.Sub(start)).String()
	},
}

func ParseFiles(path ...string) (*Template, error) {
	t, err := template.New("msg").ParseFiles(path...)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	return NewTemplate(t), nil
}

func Parse(s string) (*Template, error) {
	t, err := template.New("msg").Funcs(DefaultFuncs).Parse(s)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	return NewTemplate(t), nil
}

type Template struct {
	template *template.Template
	pool     *bytebufferpool.Pool
}

func NewTemplate(t *template.Template) *Template {
	return &Template{
		template: t,
		pool:     &bytebufferpool.Pool{},
	}
}

func (t *Template) Produce(m hook.Message) (string, error) {
	b := t.pool.Get()
	defer t.pool.Put(b)

	err := t.template.Execute(b, m)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}
