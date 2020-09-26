package models

type Task struct {
	ID            string         `json:"id,omitempty"`
	Type          string         `json:"type,omitempty"`
	Title         string         `json:"title,omitempty"`
	ResultChecker *ResultChecker `json:"result_checker,omitempty"`
}
