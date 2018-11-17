package models

type Check struct {
	Session string `json:"session, omitempty"`
	User    string `json:"user, omitempty"`
}

type Balance struct {
	Balance int `json:"balance"`
}
