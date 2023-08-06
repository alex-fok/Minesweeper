package types

type OnlineStatus struct {
	Alias    string `json:"alias"`
	IsOnline bool   `json:"isOnline"`
	IsReady  bool   `json:"isReady"`
}
