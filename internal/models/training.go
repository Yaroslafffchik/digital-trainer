package models

type Training struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	SessionID   string `json:"session_id"`
}
