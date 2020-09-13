package emitter

import (
	"context"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/tdakkota/vkalertmanager/pkg/hook"
)

type VK struct {
	client      *api.VK
	receiverIDs []int
	template    Template
	isGroup     bool
}

func NewVK(client *api.VK, receiverIDs []int, template Template) VK {
	return VK{client: client, receiverIDs: receiverIDs, template: template}
}

func (v VK) sendUser(ctxt context.Context, msg string) error {
	b := params.NewMessagesSendBuilder()
	b.Message(msg)

	for i := range v.receiverIDs {
		b.PeerID(v.receiverIDs[i])

		_, err := v.client.MessagesSend(b.Params)
		if err != nil {
			return err
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

		b.UserIDs(v.receiverIDs[i:to])
		_, err := v.client.MessagesSend(b.Params)
		if err != nil {
			return err
		}
	}

	return nil
}

func (v VK) Emit(ctxt context.Context, m hook.Message) error {
	msg, err := v.template.Produce(m)
	if err != nil {
		return err
	}

	if v.isGroup {
		err = v.sendGroup(ctxt, msg)
	} else {
		err = v.sendUser(ctxt, msg)
	}

	return err
}
