package models

type Month struct {
	Month int          `json:"month"`
	Year  int          `json:"year"`
	Items []ActiveItem `json:"items"`
}

type ActiveItem struct {
	ID          int     `json:"id"`          // one-indexed; common ID
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Accrued     int     `json:"accrued"`
	Realized    int     `json:"realized"`
	Mutable     bool    `json:"mutable"`
}

func (b *ActiveItem) Remaining() int {
	return b.Accrued - b.Realized
}