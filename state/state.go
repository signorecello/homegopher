package state

import (
	"time"
)

type State struct {
	EntityID    string     `json:"entity_id"`
	LastChanged time.Time  `json:"last_changed"`
	State       string     `json:"state"`
	Attributes  Attributes `json:"attributes"`
	LastUpdated time.Time  `json:"last_updated"`
	Context     struct {
		ID       string      `json:"id"`
		ParentID interface{} `json:"parent_id"`
		UserID   string      `json:"user_id"`
	} `json:"context"`
}

type Attributes map[string]interface{}

