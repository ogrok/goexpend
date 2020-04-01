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
	Amount		int     `json:"amount"`
	Accrued     int     `json:"accrued"`
	Realized    int     `json:"realized"`
	Mutable     bool    `json:"mutable"`
	OneTime		bool	`json:"one_time"`
}

func (b *ActiveItem) Remaining() int {
	return b.Accrued - b.Realized
}

// returns amount of excess from normal amount (manual accrual)
func (a *ActiveItem) Excess() int {
	if a.OneTime {
		return 0
	} else {
		return a.Accrued - a.Amount
	}
}