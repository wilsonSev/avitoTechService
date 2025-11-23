package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wilsonSev/avitoTechService/internal/model"
)

var ErrPRExists error = errors.New("pull request already exists")
var ErrPRNotFound error = errors.New("pull request not found")

type PRRepo struct {
	db *pgxpool.Pool
}

func NewPRRepo(db *pgxpool.Pool) *PRRepo {
	return &PRRepo{db}
}

func (r *PRRepo) Create(ctx context.Context, pr *model.PullRequest) error {
	var tmp string
	err := r.db.QueryRow(ctx,
		`SELECT id FROM pull_requests WHERE id = $1`,
		pr.PullRequestID,
	).Scan(&tmp)

	if err == nil {
		return ErrPRExists
	}
	if err != pgx.ErrNoRows {
		return err
	}
	_, err = r.db.Exec(ctx,
		`INSERT INTO pull_requests (id, name, author_id, status)
         VALUES ($1, $2, $3, 'OPEN')`,
		pr.PullRequestID, pr.PullRequestName, pr.AuthorID,
	)
	return err
}

func (r *PRRepo) Get(ctx context.Context, prID string) (*model.PullRequest, error) {
	var pr model.PullRequest
	err := r.db.QueryRow(ctx,
		`SELECT id, name, author_id, status, created_at, merged_at
         	 FROM pull_requests
         	 WHERE id = $1`,
		prID,
	).Scan(
		&pr.PullRequestID,
		&pr.PullRequestName,
		&pr.AuthorID,
		&pr.Status,
		&pr.CreatedAt,
		&pr.MergedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, ErrPRNotFound
	}
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx,
		`SELECT user_id
         	 FROM pr_reviewers
         	 WHERE pr_id = $1`,
		prID,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var reviewerID string
		if err := rows.Scan(&reviewerID); err != nil {
			return nil, err
		}
		pr.AssignedReviewers = append(pr.AssignedReviewers, reviewerID)
	}
	return &pr, nil
}

func (r *PRRepo) AddReviewers(ctx context.Context, prID string, reviewers []string) error {
	if len(reviewers) == 0 {
		return nil
	}
	for _, reviewerID := range reviewers {
		_, err := r.db.Exec(ctx,
			`INSERT INTO pr_reviewers (pr_id, user_id)
             	 VALUES ($1, $2)
        		 ON CONFLICT DO NOTHING`,
			prID, reviewerID,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *PRRepo) MergePR(ctx context.Context, prID string) error {

	var status string
	err := r.db.QueryRow(ctx,
		`SELECT status FROM pull_requests WHERE id = $1`,
		prID,
	).Scan(&status)

	if err == pgx.ErrNoRows {
		return ErrPRNotFound
	}
	if err != nil {
		return err
	}

	if status == "MERGED" {
		return nil
	}

	_, err = r.db.Exec(ctx,
		`UPDATE pull_requests
             SET status = 'MERGED', merged_at = NOW()
             WHERE id = $1`,
		prID,
	)

	return err
}

func (r *PRRepo) ReplaceReviewer(ctx context.Context, prID, oldReviewerID, newReviewerID string) error {

	var status string
	err := r.db.QueryRow(ctx,
		`SELECT status FROM pull_requests WHERE id = $1`,
		prID,
	).Scan(&status)

	if err == pgx.ErrNoRows {
		return ErrPRNotFound
	}
	if err != nil {
		return err
	}

	if status == "MERGED" {
		return errors.New("cannot replace reviewer after merge")
	}

	var tmp string
	err = r.db.QueryRow(ctx,
		`SELECT user_id
             FROM pr_reviewers
             WHERE pr_id = $1 AND user_id = $2`,
		prID, oldReviewerID,
	).Scan(&tmp)

	if err == pgx.ErrNoRows {
		return errors.New("old reviewer not assigned")
	}
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx,
		`DELETE FROM pr_reviewers
             WHERE pr_id = $1 AND user_id = $2`,
		prID, oldReviewerID,
	)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx,
		`INSERT INTO pr_reviewers (pr_id, user_id)
             VALUES ($1, $2)
             ON CONFLICT DO NOTHING`,
		prID, newReviewerID,
	)

	return err
}

func (r *PRRepo) GetPRsByReviewer(ctx context.Context, userID string) ([]model.PullRequest, error) {

	rows, err := r.db.Query(ctx,
		`SELECT pr.id, pr.name, pr.author_id, pr.status
         FROM pr_reviewers r
         JOIN pull_requests pr ON pr.id = r.pr_id
         WHERE r.user_id = $1
         ORDER BY pr.created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.PullRequest

	for rows.Next() {
		var pr model.PullRequest
		err := rows.Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status)
		if err != nil {
			return nil, err
		}
		result = append(result, pr)
	}

	return result, nil
}
