package model

type PullRequest struct {
	id        int64  `json:"pull_request_id"`
	name      string `json:"pull_request_name"`
	author    User   `json:"author_id"`
	status    string `json:"status"`
	reviewers []User `json:"assigned_reviewers"`
}
