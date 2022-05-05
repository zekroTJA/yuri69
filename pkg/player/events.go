package player

type EventType string

const (
	EventPlayStart     = EventType("playstart")
	EventPlayEnd       = EventType("playend")
	EventPlayStuck     = EventType("playstuck")
	EventPlayException = EventType("playexception")

	EventFastTrigger = EventType("fasttrigger")

	EventError = EventType("error")
)

type Event struct {
	Type    EventType
	Ident   string
	GuildID string
	UserID  string
	Err     error
}
