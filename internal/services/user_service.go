package services

import (
	"context"

	"github.com/wilsonSev/avitoTechService/internal/model"
	"github.com/wilsonSev/avitoTechService/internal/storage"
)

type UserService struct {
	userRepo *storage.UserRepo
	prRepo   *storage.PRRepo
}

func NewUserService(userRepo *storage.UserRepo, prRepo *storage.PRRepo) *UserService {
	return &UserService{userRepo, prRepo}
}

// SetIsActive /users/setIsActive
func (s *UserService) SetIsActive(ctx context.Context, userID string, active bool) (*model.User, error) {
	user, err := s.userRepo.SetActive(ctx, userID, active)
	if err == storage.ErrUserNotFound {
		return nil, err
	}
	return user, err
}

// GetUserReviews /users/getReview
func (s *UserService) GetUserReviews(ctx context.Context, userID string) ([]model.PullRequest, error) {
	exists, err := s.userRepo.Exists(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, storage.ErrUserNotFound
	}

	prs, err := s.prRepo.GetPRsByReviewer(ctx, userID)
	if err != nil {
		return nil, err
	}

	return prs, nil
}
