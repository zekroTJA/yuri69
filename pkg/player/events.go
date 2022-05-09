package player

type EventType string

const (
	EventPlayStart     = EventType("playstart")
	EventPlayEnd       = EventType("playend")
	EventPlayStuck     = EventType("playstuck")
	EventPlayException = EventType("playexception")

	EventVoiceJoin  = EventType("voicejoin")
	EventVoiceLeave = EventType("voiceleave")

	EventVoiceInit   = EventType("voiceinit")
	EventVoiceDeinit = EventType("voicedeinit")

	EventFastTrigger = EventType("fasttrigger")

	EventError = EventType("error")
)

type Event struct {
	Type    EventType `json:"type"`
	Ident   string    `json:"ident,omitempty"`
	GuildID string    `json:"guild_id,omitempty"`
	UserID  string    `json:"user_id,omitempty"`
	Err     error     `json:"error,omitempty"`
}
