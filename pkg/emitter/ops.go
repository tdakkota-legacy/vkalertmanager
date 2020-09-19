package emitter

import "github.com/tdakkota/vkalertmanager/pkg/template"

// VKOp is VK struct option function.
type VKOp func(vk *VK)

// WithIsUser sets user mode message sending.
func WithIsUser(v bool) VKOp {
	return func(vk *VK) {
		vk.isUser = v
	}
}

// WithTemplate sets message template.
func WithTemplate(t *template.Template) VKOp {
	return func(vk *VK) {
		vk.template = t
	}
}
