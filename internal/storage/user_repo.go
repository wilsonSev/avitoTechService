package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wilsonSev/avitoTechService/internal/model"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db}
}

func (r *UserRepo) GetByID(ctx context.Context, userID string) (*model.User, error) {
	var u model.User

	err := r.db.QueryRow(ctx,
		`SELECT id, username, team_name, is_active
         	 FROM users
         	 WHERE id = $1`,
		userID,
	).Scan(&u.UserID, &u.Username, &u.TeamName, &u.IsActive)

	if err == pgx.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) SetActive(ctx context.Context, userID string, isActive bool) (*model.User, error) {
	var u model.User

	err := r.db.QueryRow(ctx,
		`UPDATE users
         SET is_active = $1
         WHERE id = $2
         RETURNING id, username, team_name, is_active`,
		isActive, userID,
	).Scan(&u.UserID, &u.Username, &u.TeamName, &u.IsActive)

	if err == pgx.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserRepo) GetTeamUsers(ctx context.Context, teamName string) ([]model.User, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, username, team_name, is_active
         	 FROM users
         	 WHERE team_name = $1`,
		teamName,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User

	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.UserID, &u.Username, &u.TeamName, &u.IsActive); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepo) Exists(ctx context.Context, userID string) (bool, error) {
	var dummy string
	err := r.db.QueryRow(ctx,
		`SELECT id FROM users WHERE id = $1`,
		userID,
	).Scan(&dummy)

	if err == pgx.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
