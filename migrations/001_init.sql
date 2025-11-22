CREATE TABLE users (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    team_name TEXT NOT NULL
    is_active BOOLEAN
);

CREATE TABLE pull_request (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    author_id TEXT NOT NULL REFERENCES users(id),
    status TEXT NOT NULL CHECK (status IN ('OPEN', 'MERGED')),
    reviewers TEXT[] NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ,
    merged_at TIMESTAMPTZ,
);

CREATE TABLE team (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    reviewers TEXT[] NOT NULL REFERENCES users(id),
)