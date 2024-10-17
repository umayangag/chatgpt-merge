package models

type Snippet struct {
	CreateTime float64     `json:"timestamp"`
	Content    interface{} `json:"content"`
	Role       string      `json:"role"`
}
