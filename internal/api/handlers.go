package api

import (
	"encoding/json"
	"github.com/wilsonSev/avitoTechService/internal/services"
	"net/http"

	"github.com/wilsonSev/avitoTechService/internal/model"
)

type Handlers struct {
	prService   *services.PRService
	teamService *services.TeamService
	userService *services.UserService
}

func NewHandlers(
	pr *services.PRService,
	team *services.TeamService,
	user *services.UserService,
) *Handlers {
	return &Handlers{
		prService:   pr,
		teamService: team,
		userService: user,
	}
}

func (Handlers) GetTeam() {

}

func (Handlers) SetUserActive() {

}

func (Handlers) GetReviewedPRs() {

}

func (Handlers) CreatePR() {

}

func (Handlers) MergePR() {

}

func (Handlers) Reassign() {

}
