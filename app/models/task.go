package models

type Task struct {
	ID            string        `json:"id"`
	Type          string        `json:"type"`
	Title         string        `json:"title"`
	ResultChecker ResultChecker `json:"result_checker"`
}
