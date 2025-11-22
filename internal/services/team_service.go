package services

import "github.com/wilsonSev/avitoTechService/internal/storage"

type TeamService struct {
	TeamRepo *storage.TeamRepo
}

func NewTeamService(teamRepo *storage.TeamRepo) *TeamService {
	return &TeamService{teamRepo}
}
