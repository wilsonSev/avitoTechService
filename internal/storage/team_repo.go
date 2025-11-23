package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wilsonSev/avitoTechService/internal/model"
)

var ErrTeamExists = errors.New("team already exists")

type TeamRepo struct {
	db *pgxpool.Pool
}

func NewTeamRepo(db *pgxpool.Pool) *TeamRepo {
	return &TeamRepo{db}
}

// Создать команду и всех ее участников
func (r *TeamRepo) CreateTeamWithMembers(ctx context.Context, teamName string, members []model.TeamMember) error {

	// создание новой транзакции
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var tmp string

	err = tx.QueryRow(ctx,
		`SELECT name FROM teams WHERE name = $1`,
		teamName,
	).Scan(&tmp)

	if err == nil {
		return ErrTeamExists
	}

	if err != pgx.ErrNoRows {
		return err
	}

	// создаем команду по имени
	_, err = tx.Exec(ctx,
		`INSERT INTO teams (name) VALUES ($1)`,
		teamName,
	)
	if err != nil {
		return err
	}

	// Добавляем пользователей
	for _, m := range members {
		_, err = tx.Exec(ctx,
			`INSERT INTO users (id, username, is_active, team_name)
				 VALUES ($1, $2, $3, $4)
				 ON CONFLICT (id) DO UPDATE
				 	SET username = EXCLUDED.username,
                     	is_active = EXCLUDED.is_active,
                     	team_name = EXCLUDED.team_name`,
			m.UserID, m.Username, m.IsActive, teamName,
		)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *TeamRepo) GetTeam(ctx context.Context, name string) (*model.Team, error) {
	var team model.Team
	team.TeamName = name

	rows, err := r.db.Query(ctx,
		`SELECT id, username, is_active
			 FROM users WHERE team_name = $1`, name)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var member model.TeamMember
		err = rows.Scan(&member.UserID, &member.Username, &member.IsActive)
		if err != nil {
			return nil, err
		}
		team.Members = append(team.Members, member)
	}
	if len(team.Members) == 0 {
		return nil, pgx.ErrNoRows
	}
	return &team, nil
}
