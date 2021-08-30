package events

import (
	"time"
	"github.com/signorecello/homegopher/state"
)

type StateChanged struct {
	EntityID string `json:"entity_id"`
	NewState state.State  `json:"new_state"`
	OldState state.State  `json:"old_state"`
}


type StateChangedEvent struct {
	ID    int    `json:"id"`
	Type  string `json:"type"`
	Event struct {
		EventType string       `json:"event_type"`
		Data      StateChanged `json:"data"`
		TimeFired time.Time    `json:"time_fired"`
		Origin    string       `json:"origin"`
		Context   struct {
			ID       string      `json:"id"`
			ParentID interface{} `json:"parent_id"`
			UserID   string      `json:"user_id"`
		} `json:"context"`
	} `json:"event"`
}
