package services

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/wilsonSev/avitoTechService/internal/model"
	"github.com/wilsonSev/avitoTechService/internal/storage"
)

type PRService struct {
	prRepo   *storage.PRRepo
	userRepo *storage.UserRepo
}

func NewPRService(prRepo *storage.PRRepo, userRepo *storage.UserRepo) *PRService {
	return &PRService{prRepo, userRepo}
}

// CreatePR /pullRequest/create
func (s *PRService) CreatePR(ctx context.Context, req model.PullRequest) (*model.PullRequest, error) {
	author, err := s.userRepo.GetByID(ctx, req.AuthorID)
	if err == storage.ErrUserNotFound {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	teamUsers, err := s.userRepo.GetTeamUsers(ctx, author.TeamName)
	if err != nil {
		return nil, err
	}

	var candidates []model.User
	for _, u := range teamUsers {
		if u.IsActive && u.UserID != author.UserID {
			candidates = append(candidates, u)
		}
	}

	rand.Seed(time.Now().UnixNano())

	var reviewers []string

	if len(candidates) == 1 {
		reviewers = []string{candidates[0].UserID}
	} else if len(candidates) >= 2 {
		rand.Shuffle(len(candidates), func(i, j int) {
			candidates[i], candidates[j] = candidates[j], candidates[i]
		})
		reviewers = []string{candidates[0].UserID, candidates[1].UserID}
	}

	err = s.prRepo.Create(ctx, &req)
	if err == storage.ErrPRExists {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	err = s.prRepo.AddReviewers(ctx, req.PullRequestID, reviewers)
	if err != nil {
		return nil, err
	}

	created, err := s.prRepo.Get(ctx, req.PullRequestID)
	if err != nil {
		return nil, err
	}

	return created, nil
}

// MergePR /pullRequest/merge
func (s *PRService) MergePR(ctx context.Context, prID string) (*model.PullRequest, error) {
	err := s.prRepo.MergePR(ctx, prID)
	if err != nil {
		return nil, err
	}

	pr, err := s.prRepo.Get(ctx, prID)
	if err != nil {
		return nil, err
	}

	return pr, nil
}

// ReassignReviewer pullRequest/reassign
func (s *PRService) ReassignReviewer(ctx context.Context, prID, oldReviewerID string) (*model.PullRequest, string, error) {
	pr, err := s.prRepo.Get(ctx, prID)
	if err == storage.ErrPRNotFound {
		return nil, "", err
	}
	if err != nil {
		return nil, "", err
	}

	if pr.Status == "MERGED" {
		return nil, "", errors.New("cannot reassign on merged PR")
	}

	found := false
	for _, r := range pr.AssignedReviewers {
		if r == oldReviewerID {
			found = true
			break
		}
	}
	if !found {
		return nil, "", errors.New("reviewer is not assigned")
	}

	oldReviewer, err := s.userRepo.GetByID(ctx, oldReviewerID)
	if err != nil {
		return nil, "", err
	}

	teamUsers, err := s.userRepo.GetTeamUsers(ctx, oldReviewer.TeamName)
	if err != nil {
		return nil, "", err
	}

	var candidates []model.User
	assigned := map[string]bool{}
	for _, r := range pr.AssignedReviewers {
		assigned[r] = true
	}

	for _, u := range teamUsers {
		if u.IsActive && !assigned[u.UserID] {
			candidates = append(candidates, u)
		}
	}

	if len(candidates) == 0 {
		return nil, "", errors.New("no active replacement candidate in team")
	}

	rand.Seed(time.Now().UnixNano())
	newReviewer := candidates[rand.Intn(len(candidates))].UserID

	err = s.prRepo.ReplaceReviewer(ctx, prID, oldReviewerID, newReviewer)
	if err != nil {
		return nil, "", err
	}

	updated, err := s.prRepo.Get(ctx, prID)
	if err != nil {
		return nil, "", err
	}

	return updated, newReviewer, nil
}
