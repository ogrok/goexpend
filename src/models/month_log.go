package models

type MonthLog struct {
	ID       int            `json:"id"`        // iterative ID just to keep the list ordinal
	Month    int            `json:"month"`
	Year     int            `json:"year"`
	LogItems []MonthLogItem `json:"log_items"`
}

type MonthLogItem struct {
	Name       string  `json:"name"`
	Category   string  `json:"category"`
	Accrued    float64 `json:"accrued"`
	Realized   float64 `json:"realized"`
	Mutable    bool    `json:"mutable"`
}