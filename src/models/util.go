package models

type Modification struct {
	ID          int
	Amount      int
	Category    string
	Description string
	Name        string
	Realized    int
}

type ActiveItemView struct {
	ID              int
	Name            string
	Category        string
	Description     string
	CurrentAccrued  int
	Realized        int
	Mutable         bool
	Amount          int
	Recurrence      string
	RecurrenceMonth int
}

func (v *ActiveItemView) Remains() int {
	return v.CurrentAccrued - v.Realized
}