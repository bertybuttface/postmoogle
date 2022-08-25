package bot

import (
	"context"

	"github.com/getsentry/sentry-go"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

func (b *Bot) initSync() {
	b.lp.OnEventType(
		event.StateMember,
		func(_ mautrix.EventSource, evt *event.Event) {
			// Trying to debug the membership=join event being handled twice here.
			eventJSON, _ := evt.MarshalJSON()
			b.log.Debug(string(eventJSON))

			go b.onMembership(evt)
		},
	)
	b.lp.OnEventType(
		event.EventMessage,
		func(_ mautrix.EventSource, evt *event.Event) {
			go b.onMessage(evt)
		})
	b.lp.OnEventType(
		event.EventEncrypted,
		func(_ mautrix.EventSource, evt *event.Event) {
			go b.onEncryptedMessage(evt)
		})
}

func (b *Bot) onMembership(evt *event.Event) {
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetUser(sentry.User{ID: evt.Sender.String()})
		scope.SetContext("event", map[string]string{
			"id":     evt.ID.String(),
			"room":   evt.RoomID.String(),
			"sender": evt.Sender.String(),
		})
	})

	if evt.Sender == b.lp.GetClient().UserID {
		// Handle membership events related to our own (bot) user first

		switch evt.Content.AsMember().Membership {
		case event.MembershipJoin:
			b.onBotJoin(evt, hub)
		}

		return
	}

	// Handle membership events related to other users
}

func (b *Bot) onMessage(evt *event.Event) {
	// ignore own messages
	if evt.Sender == b.lp.GetClient().UserID {
		return
	}

	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetUser(sentry.User{ID: evt.Sender.String()})
		scope.SetContext("event", map[string]string{
			"id":     evt.ID.String(),
			"room":   evt.RoomID.String(),
			"sender": evt.Sender.String(),
		})
	})
	ctx := sentry.SetHubOnContext(context.Background(), hub)
	span := sentry.StartSpan(ctx, "http.server", sentry.TransactionName("onMessage"))
	defer span.Finish()

	b.handle(span.Context(), evt)
}

func (b *Bot) onEncryptedMessage(evt *event.Event) {
	// ignore own messages
	if evt.Sender == b.lp.GetClient().UserID {
		return
	}

	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetUser(sentry.User{ID: evt.Sender.String()})
		scope.SetContext("event", map[string]string{
			"id":     evt.ID.String(),
			"room":   evt.RoomID.String(),
			"sender": evt.Sender.String(),
		})
	})
	ctx := sentry.SetHubOnContext(context.Background(), hub)
	span := sentry.StartSpan(ctx, "http.server", sentry.TransactionName("onMessage"))
	defer span.Finish()

	decrypted, err := b.lp.GetMachine().DecryptMegolmEvent(evt)
	if err != nil {
		b.Error(span.Context(), evt.RoomID, "cannot decrypt a message: %v", err)
		return
	}

	b.handle(span.Context(), decrypted)
}

// onBotJoin handles the "bot joined the room" event
func (b *Bot) onBotJoin(evt *event.Event, hub *sentry.Hub) {
	ctx := sentry.SetHubOnContext(context.Background(), hub)
	span := sentry.StartSpan(ctx, "http.server", sentry.TransactionName("onBotJoin"))
	defer span.Finish()

	b.sendIntroduction(ctx, evt.RoomID)
	b.sendHelp(ctx, evt.RoomID)
}
