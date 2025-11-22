package services

import (
	"github.com/wilsonSev/avitoTechService/internal/storage"
)

type PRService struct {
	prRepo *storage.PRRepo
}

func NewUserService(prRepo *storage.PRRepo) *UserService {
	return &UserService{}
}
