package models

type ItemTemplate struct {
	ID         int     `json:"id"`            // one-indexed
	Name       string  `json:"name"`
	Category   string  `json:"category"`
	Amount     float64 `json:"amount"`
	Recurrence string  `json:"recurrence"`
	Mutable    bool    `json:"mutable"`
}