package store

type Event struct {
	Id          int
	AggregateID string
	EventData   string
	EventType   int
}
