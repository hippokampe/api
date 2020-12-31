package models

type Projects []Project

type Project struct {
	exists   bool
	ID       string     `json:"id"`
	Title    string     `json:"title"`
	Category string     `json:"category"`
	Score    string     `json:"score"`
	Tasks    TasksBasic `json:"tasks"`
}

func (p Project) Exists() bool {
	return p.Title != ""
}
