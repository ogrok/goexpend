package models

type Month struct {
	Month int         `json:"month"`
	Year  int         `json:"year"`
	Items []MonthItem `json:"items"`
}

type MonthItem struct {
	ID         int     `json:"id"`          // one-indexed; common ID
	Name       string  `json:"name"`
	Category   string  `json:"category"`
	InstanceOf int     `json:"instance_of"` // id of parent ItemTemplate
	Accrued    float64 `json:"accrued"`
	Realized   float64 `json:"realized"`
	Mutable    bool    `json:"mutable"`
}

func (b *MonthItem) Remaining() float64 {
	return b.Accrued - b.Realized
}