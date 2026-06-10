package events

type VersionedEvent interface {
	Event
	Version() int
}

type VersionedBaseEvent struct {
	BaseEvent
	V int `json:"version"`
}

func (e VersionedBaseEvent) Version() int {
	return e.V
}

const (
	EventVersionV1 = 1
	EventVersionV2 = 2
)
