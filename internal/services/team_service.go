package services

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/wilsonSev/avitoTechService/internal/model"
	"github.com/wilsonSev/avitoTechService/internal/storage"
)

type TeamService struct {
	teamRepo *storage.TeamRepo
	userRepo *storage.UserRepo
}

func NewTeamService(teamRepo *storage.TeamRepo, userRepo *storage.UserRepo) *TeamService {
	return &TeamService{teamRepo, userRepo}
}

// CreateTeam handles /team/add
func (s *TeamService) CreateTeam(ctx context.Context, team model.Team) (*model.Team, error) {
	err := s.teamRepo.CreateTeamWithMembers(ctx, team.TeamName, team.Members)
	if err == storage.ErrTeamExists {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &team, nil
}

// GetTeam handles /team/get
func (s *TeamService) GetTeam(ctx context.Context, name string) (*model.Team, error) {
	team, err := s.teamRepo.GetTeam(ctx, name)
	if err == pgx.ErrNoRows {
		return nil, err
	}
	return team, err
}
