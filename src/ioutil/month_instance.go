package ioutil

type Month struct {
	Month int         `json:"month"`
	Year  int         `json:"year"`
	Items []MonthItem `json:"items"`
}

type MonthItem struct {
	ID          int     `json:"id"`          // one-indexed; common ID
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Accrued     int     `json:"accrued"`
	Realized    int     `json:"realized"`
	Mutable     bool    `json:"mutable"`
}

func (b *MonthItem) Remaining() int {
	return b.Accrued - b.Realized
}

// returns amount of excess from template accrued amount, e.g. amount of manual accrual at the moment
func (b *MonthItem) Excess() int {
	template, err := GetSpecificTemplate(b.ID)

	if err != nil {
		return 0
	}

	return b.Accrued - template.Amount
}