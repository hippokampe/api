package models

type Task struct {
	ID              string         `json:"id,omitempty"`
	Type            string         `json:"type,omitempty"`
	Title           string         `json:"title,omitempty"`
	Done            bool           `json:"done,omitempty"`
	FileDescription string         `json:"file_description"`
	ResultChecker   *ResultChecker `json:"result_checker,omitempty"`
}
