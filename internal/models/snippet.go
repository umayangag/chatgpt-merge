package models

type Snippet struct {
	CreateTime float64 `json:"timestamp"`
	TimeString string  `json:"timestring"`
	Role       string  `json:"role"`
	Content    string  `json:"content"`
}
