package ioutil

type MonthLog struct {
	ID       int            `json:"id"`        // iterative ID just to keep the list ordinal
	Month    int            `json:"month"`
	Year     int            `json:"year"`
	LogItems []MonthLogItem `json:"log_items"`
}

type MonthLogItem struct {
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Accrued     int     `json:"accrued"`
	Excess		int     `json:"excess"` 	   // excess accrued
	Realized    int     `json:"realized"`
	Remaining   int     `json:"remaining"`     // accrued - realized
	Mutable     bool    `json:"mutable"`
}