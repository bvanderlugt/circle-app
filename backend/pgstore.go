package main

import (
	"context"
	"database/sql"

	"gocloud.dev/postgres"
)

type PgStore struct {
	db *sql.DB
}

func NewPgStore(ctx context.Context, url string) (*PgStore, error) {
	db, err := postgres.Open(ctx, url)
	if err != nil {
		return nil, err
	}
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}
	if err := ensureSchema(ctx, db); err != nil {
		db.Close()
		return nil, err
	}
	return &PgStore{db: db}, nil
}

func (p *PgStore) Close() error { return p.db.Close() }

func ensureSchema(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS circles (
			id         SERIAL PRIMARY KEY,
			name       TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)
	return err
}

func (p *PgStore) Count(ctx context.Context) (int, error) {
	var n int
	err := p.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM circles`).Scan(&n)
	return n, err
}

func (p *PgStore) List(ctx context.Context) ([]Circle, error) {
	rows, err := p.db.QueryContext(ctx, `SELECT id, name, created_at FROM circles ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []Circle{}
	for rows.Next() {
		var c Circle
		if err := rows.Scan(&c.ID, &c.Name, &c.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (p *PgStore) Add(ctx context.Context, name string) (Circle, error) {
	c := Circle{Name: name}
	err := p.db.QueryRowContext(ctx,
		`INSERT INTO circles (name) VALUES ($1) RETURNING id, created_at`,
		name,
	).Scan(&c.ID, &c.CreatedAt)
	return c, err
}

func (p *PgStore) Delete(ctx context.Context, id int64) error {
	res, err := p.db.ExecContext(ctx, `DELETE FROM circles WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
