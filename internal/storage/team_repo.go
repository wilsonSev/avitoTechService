package storage

import "github.com/jackc/pgx/v5/pgxpool"

type TeamRepo struct {
	db *pgxpool.Pool
}

func NewTeamRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db}
}
