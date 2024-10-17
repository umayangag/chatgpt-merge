package models

type Conversation struct {
	Mapping map[string]struct {
		Message    Message `json:"message"`
		CreateTime string  `json:"create_time"`
	} `json:"mapping"`
}

type Message struct {
	CreateTime float64 `json:"create_time"`
	Content    Content `json:"content"`
	Author     Author  `json:"author"`
}

type Content struct {
	Parts interface{} `json:"parts"`
}

type Author struct {
	Role string `json:"role"`
}
