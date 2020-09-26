package models

type ResultChecker struct {
	Checks []Check `json:"checks,omitempty"`
}

type Check struct {
	ID       string `json:"id,omitempty"`
	Type     string `json:"type,omitempty"`
	Status   bool   `json:"status,omitempty"`
	Position int    `json:"position"`
}
