package models

type Metrics struct {
	ID        int     `json:"id"`
	SessionID string  `json:"session_id"`
	Pulse     int     `json:"pulse"`
	Speed     float64 `json:"speed"`
	Timestamp string  `json:"timestamp"`
}
