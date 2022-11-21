package bot

import (
	"context"
	"strings"

	"gitlab.com/etke.cc/postmoogle/utils"
)

func (b *Bot) handle(ctx context.Context) {
	evt := eventFromContext(ctx)
	err := b.lp.GetClient().MarkRead(evt.RoomID, evt.ID)
	if err != nil {
		b.log.Error("cannot send read receipt: %v", err)
	}

	content := evt.Content.AsMessage()
	if content == nil {
		b.Error(ctx, evt.RoomID, "cannot read message")
		return
	}
	message := strings.TrimSpace(content.Body)
	cmd := b.parseCommand(message, true)
	if cmd == nil {
		if utils.EventParent("", content) != "" {
			b.SendEmailReply(ctx)
		}
		return
	}

	b.handleCommand(ctx, evt, cmd)
}
