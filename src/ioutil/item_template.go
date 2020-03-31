package ioutil

type ItemTemplate struct {
	ID              int     `json:"id"`               // one-indexed; common ID
	Name            string  `json:"name"`
	Category        string  `json:"category"`
	Description     string  `json:"description"`
	Amount          int     `json:"amount"`
	Recurrence      string  `json:"recurrence"`
	RecurrenceMonth int		`json:"recurrence_month"` // 0 if recurrence != yearly, else 1-12
	Mutable         bool    `json:"mutable"`
}