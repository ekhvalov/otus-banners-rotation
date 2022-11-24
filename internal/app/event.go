package app

type EventType string

const (
	EventClick  EventType = "click"
	EventSelect EventType = "select"
)

type Event struct {
	Type           EventType
	SlotID         string
	BannerID       string
	SocialGroupID  string
	TimestampMicro int64
}
