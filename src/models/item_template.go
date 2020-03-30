package models

type ItemTemplate struct {
	ID              int     `json:"id"`               // one-indexed; user does not interact with this
	Name            string  `json:"name"`
	Category        string  `json:"category"`
	Amount          float64 `json:"amount"`
	Recurrence      string  `json:"recurrence"`
	RecurrenceMonth int		`json:"recurrence_month"` // 0 if recurrence != yearly, else 1-12
	Mutable         bool    `json:"mutable"`
}