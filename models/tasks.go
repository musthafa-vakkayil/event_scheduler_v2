package models

type LogPayload struct {
	LogID int `json:"log_id"`
}

type EventPayload struct {
	EventID int64 `json:"event_id"`
}
