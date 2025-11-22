package model

type User struct {
	id   int64  `json:"user_id"`
	name string `json:"username"`
	// team     string `json:"team_name"`
	isActive bool `json:"is_active"`
}
