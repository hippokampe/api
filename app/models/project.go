package models

type Project struct {
	Category            string          `json:"category,omitempty"`
	Title               string          `json:"title,omitempty"`
	ID                  string          `json:"code,omitempty"`
	Score               string          `json:"score,omitempty"`
	Tasks               []Task          `json:"tasks,omitempty"`
	DeadlineInformation DeadlineSummary `json:"deadline_information,omitempty"`
}

type Projects struct {
	CurrentProjects CurrentProjects `json:"current_projects"`
	AllProjects     []Project       `json:"all_projects"`
}

type CurrentProjects struct {
	FirstDeadline  []Project `json:"first_deadline"`
	SecondDeadline []Project `json:"second_deadline"`
	Total          int       `json:"total"`
}

type DeadlineSummary struct {
	Period        int    `json:"period,omitempty"`
	Started       string `json:"started,omitempty"`
	Finished      string `json:"finished,omitempty"`
	RemainingDate int    `json:"remaining_date,omitempty"`
}
