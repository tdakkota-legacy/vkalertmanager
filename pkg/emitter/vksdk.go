package emitter

import (
	"context"
	"fmt"

	"github.com/tdakkota/vkalertmanager/pkg/template"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"

	"github.com/tdakkota/vkalertmanager/pkg/hook"
)

// VK is VK API message emitter.
type VK struct {
	client      *api.VK
	receiverIDs []int
	template    *template.Template
	isUser      bool
}

// NewVK creates new VK struct.
func NewVK(client *api.VK, receiverIDs []int, ops ...VKOp) VK {
	vk := VK{client: client, receiverIDs: receiverIDs}

	for _, op := range ops {
		op(&vk)
	}

	if vk.template == nil {
		vk.template = template.Default()
	}

	return vk
}

func (v VK) sendUser(ctxt context.Context, msg string) error {
	b := params.NewMessagesSendBuilder()
	b.Message(msg)

	for _, userID := range v.receiverIDs {
		b.PeerID(userID)

		_, err := v.client.MessagesSend(b.Params)
		if err != nil {
			return fmt.Errorf("send message to user %d: %w", userID, err)
		}
	}

	return nil
}

func (v VK) sendGroup(ctxt context.Context, msg string) error {
	b := params.NewMessagesSendBuilder()
	b.Message(msg)

	length := len(v.receiverIDs)
	for i := 0; i < length; i += 100 {
		to := i + 100
		if to > length {
			to = length
		}

		b.PeerIDs(v.receiverIDs[i:to])
		_, err := v.client.MessagesSend(b.Params)
		if err != nil {
			return fmt.Errorf("send message: %w", err)
		}
	}

	return nil
}

func (v VK) Emit(ctxt context.Context, m hook.Message) error {
	msg, err := v.template.Produce(m)
	if err != nil {
		return err
	}

	if v.isUser {
		err = v.sendUser(ctxt, msg)
	} else {
		err = v.sendGroup(ctxt, msg)
	}

	return err
}
