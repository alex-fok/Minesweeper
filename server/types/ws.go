package types

type ClientId string

type ClientMeta struct {
	Id       ClientId
	Alias    string
	IsOnline bool
	IsReady  bool
}

type Action struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}
