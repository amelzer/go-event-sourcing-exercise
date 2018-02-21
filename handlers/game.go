package handlers

import (
	"github.com/wwgberlin/go-event-sourcing-exercise/chess"
	"github.com/wwgberlin/go-event-sourcing-exercise/store"
)

const (
	EventNone int = iota
	EventMoveRequest
	EventMoveSuccess
	EventMoveFail
	EventPromotionRequest
	EventPromotionSuccess
	EventPromotionFail
	EventWhiteWins
	EventBlackWins
	EventDraw
)

func BuildGame(eventStore *store.EventStore, gameID string, lastEventID int) *chess.Game {
	game := chess.NewGame()

	for _, event := range eventStore.GetEvents() {
		if event.AggregateID != gameID {
			continue
		}

		switch event.EventType {
		case EventMoveSuccess:
			game.Move(chess.ParseMove(event.EventData))
		case EventPromotionSuccess:
			game.Promote(chess.ParsePromotion(event.EventData))
		}
	}
	return game
}

func MoveHandler(eventStore *store.EventStore, event store.Event) {
	if event.EventType != EventMoveRequest {
		return
	}
	game := BuildGame(eventStore, event.AggregateID, -1)

	ev := store.Event{
		AggregateID: event.AggregateID,
	}
	if err := game.Move(chess.ParseMove(event.EventData)); err != nil {
		ev.EventType = EventMoveFail
		ev.EventData = err.Error()
	} else {
		ev.EventType = EventMoveSuccess
		ev.EventData = event.EventData
	}
	eventStore.Persist(ev)
}

func PromotionHandler(eventStore *store.EventStore, event store.Event) {
	if event.EventType != EventPromotionRequest {
		return
	}
	game := BuildGame(eventStore, event.AggregateID, -1)

	ev := store.Event{
		AggregateID: event.AggregateID,
	}
	if err := game.Promote(chess.ParsePromotion(event.EventData)); err != nil {
		ev.EventType = EventPromotionFail
		ev.EventData = err.Error()
	} else {
		ev.EventType = EventPromotionSuccess
		ev.EventData = event.EventData
	}
	eventStore.Persist(ev)
}

func StatusChangeHandler(eventStore *store.EventStore, event store.Event) {
	game := BuildGame(eventStore, event.AggregateID, -1)
	status := game.Status()
	if status == 0 {
		return
	}

	ev := store.Event{
		AggregateID: event.AggregateID,
		EventData:   event.EventData,
	}
	if status == 1 {
		ev.EventType = EventWhiteWins
	} else if status == 2 {
		ev.EventType = EventBlackWins
	} else if status == 3 {
		ev.EventType = EventDraw
	}
	eventStore.Persist(ev)
}