package config

import (
	"strings"

	"gitlab.com/etke.cc/postmoogle/email"
	"gitlab.com/etke.cc/postmoogle/utils"
)

// account data key
const acRoomKey = "cc.etke.postmoogle.settings"

type Room map[string]string

// option keys
const (
	RoomActive        = ".active"
	RoomOwner         = "owner"
	RoomMailbox       = "mailbox"
	RoomDomain        = "domain"
	RoomNoSend        = "nosend"
	RoomNoReplies     = "noreplies"
	RoomNoCC          = "nocc"
	RoomNoSender      = "nosender"
	RoomNoRecipient   = "norecipient"
	RoomNoSubject     = "nosubject"
	RoomNoHTML        = "nohtml"
	RoomNoThreads     = "nothreads"
	RoomNoFiles       = "nofiles"
	RoomNoInlines     = "noinlines"
	RoomPassword      = "password"
	RoomSpamcheckDKIM = "spamcheck:dkim"
	RoomSpamcheckSMTP = "spamcheck:smtp"
	RoomSpamcheckSPF  = "spamcheck:spf"
	RoomSpamcheckMX   = "spamcheck:mx"
	RoomSpamlist      = "spamlist"
)

// Get option
func (s Room) Get(key string) string {
	return s[strings.ToLower(strings.TrimSpace(key))]
}

// Set option
func (s Room) Set(key, value string) {
	s[strings.ToLower(strings.TrimSpace(key))] = value
}

func (s Room) Mailbox() string {
	return s.Get(RoomMailbox)
}

func (s Room) Domain() string {
	return s.Get(RoomDomain)
}

func (s Room) Owner() string {
	return s.Get(RoomOwner)
}

func (s Room) Active() bool {
	return utils.Bool(s.Get(RoomActive))
}

func (s Room) Password() string {
	return s.Get(RoomPassword)
}

func (s Room) NoSend() bool {
	return utils.Bool(s.Get(RoomNoSend))
}

func (s Room) NoReplies() bool {
	return utils.Bool(s.Get(RoomNoReplies))
}

func (s Room) NoCC() bool {
	return utils.Bool(s.Get(RoomNoCC))
}

func (s Room) NoSender() bool {
	return utils.Bool(s.Get(RoomNoSender))
}

func (s Room) NoRecipient() bool {
	return utils.Bool(s.Get(RoomNoRecipient))
}

func (s Room) NoSubject() bool {
	return utils.Bool(s.Get(RoomNoSubject))
}

func (s Room) NoHTML() bool {
	return utils.Bool(s.Get(RoomNoHTML))
}

func (s Room) NoThreads() bool {
	return utils.Bool(s.Get(RoomNoThreads))
}

func (s Room) NoFiles() bool {
	return utils.Bool(s.Get(RoomNoFiles))
}

func (s Room) NoInlines() bool {
	return utils.Bool(s.Get(RoomNoInlines))
}

func (s Room) SpamcheckDKIM() bool {
	return utils.Bool(s.Get(RoomSpamcheckDKIM))
}

func (s Room) SpamcheckSMTP() bool {
	return utils.Bool(s.Get(RoomSpamcheckSMTP))
}

func (s Room) SpamcheckSPF() bool {
	return utils.Bool(s.Get(RoomSpamcheckSPF))
}

func (s Room) SpamcheckMX() bool {
	return utils.Bool(s.Get(RoomSpamcheckMX))
}

func (s Room) Spamlist() []string {
	return utils.StringSlice(s.Get(RoomSpamlist))
}

func (s Room) MigrateSpamlistSettings() {
	uniq := map[string]struct{}{}
	emails := utils.StringSlice(s.Get("spamlist:emails"))
	localparts := utils.StringSlice(s.Get("spamlist:localparts"))
	hosts := utils.StringSlice(s.Get("spamlist:hosts"))
	list := utils.StringSlice(s.Get(RoomSpamlist))
	delete(s, "spamlist:emails")
	delete(s, "spamlist:localparts")
	delete(s, "spamlist:hosts")

	for _, email := range emails {
		if email == "" {
			continue
		}
		uniq[email] = struct{}{}
	}

	for _, localpart := range localparts {
		if localpart == "" {
			continue
		}
		uniq[localpart+"@*"] = struct{}{}
	}

	for _, host := range hosts {
		if host == "" {
			continue
		}
		uniq["*@"+host] = struct{}{}
	}

	for _, item := range list {
		if item == "" {
			continue
		}
		uniq[item] = struct{}{}
	}

	spamlist := make([]string, 0, len(uniq))
	for item := range uniq {
		spamlist = append(spamlist, item)
	}
	s.Set(RoomSpamlist, strings.Join(spamlist, ","))
}

// ContentOptions converts room display settings to content options
func (s Room) ContentOptions() *email.ContentOptions {
	return &email.ContentOptions{
		CC:        !s.NoCC(),
		HTML:      !s.NoHTML(),
		Sender:    !s.NoSender(),
		Recipient: !s.NoRecipient(),
		Subject:   !s.NoSubject(),
		Threads:   !s.NoThreads(),

		ToKey:         "cc.etke.postmoogle.to",
		CcKey:         "cc.etke.postmoogle.cc",
		FromKey:       "cc.etke.postmoogle.from",
		RcptToKey:     "cc.etke.postmoogle.rcptTo",
		SubjectKey:    "cc.etke.postmoogle.subject",
		InReplyToKey:  "cc.etke.postmoogle.inReplyTo",
		MessageIDKey:  "cc.etke.postmoogle.messageID",
		ReferencesKey: "cc.etke.postmoogle.references",
	}
}
