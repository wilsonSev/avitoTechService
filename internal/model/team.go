package model

type Team struct {
	name    string `json:"team_name"`
	members []User `json:"members"`
}
