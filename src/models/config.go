package models

type Config struct {
	CurrentMonth  int `json:"current_month"`
	CurrentYear   int `json:"current_year"`

	// for determining when to auto-generate new MonthItems
	MonthEnd      int `json:"month_end"`

	// for checking when to ask about turn-of-month
	AskAgainAfter int `json:"ask_again_after"`
}