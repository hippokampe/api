package models

type Projects []Project
type TasksBasic []TaskBasic

type Project struct {
	exists   bool
	ID       string     `json:"id"`
	Title    string     `json:"title"`
	Category string     `json:"category"`
	Score    string     `json:"score"`
	Tasks    TasksBasic `json:"tasks"`
}

type TaskBasic struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Position string `json:"position"`
	Type     string `json:"type"`
}

func (p Project) Exists() bool {
	return p.Title != ""
}
