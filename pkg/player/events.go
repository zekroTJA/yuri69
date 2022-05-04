package player

type EventType string

const (
	EventPlayStart     = EventType("playstart")
	EventPlayEnd       = EventType("playend")
	EventPlayStuck     = EventType("playstuck")
	EventPlayException = EventType("playexception")
)

type Event struct {
	Type    EventType
	Ident   string
	GuildID string
	Err     error
}
