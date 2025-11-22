package storage

import "github.com/jackc/pgx/v5/pgxpool"

type PRRepo struct {
	db *pgxpool.Pool
}

func NewPRRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db}
}
