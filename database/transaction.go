package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBTX interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
}

func NewDBTX(dsn string) (DBTX, error) {
	ctx := context.Background()

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}

type Transaction struct {
	dbtx DBTX
}

func NewTransaction(dbtx DBTX) (Transaction, error) {
	return Transaction{dbtx: dbtx}, nil
}

func (tx Transaction) InsertSample(ctx context.Context, name string) (Sample, error) {
	var id string

	if err := tx.dbtx.QueryRow(ctx, `INSERT INTO samples (name) VALUES ($1) RETURNING id`, name).Scan(&id); err != nil {
		return Sample{}, err
	}

	return Sample{ID: id, Name: name}, nil
}

func (tx Transaction) UpdateSample(ctx context.Context, id string, name string) error {
	if _, err := tx.dbtx.Exec(ctx, `UPDATE samples SET name = $1 WHERE id = $2`, name, id); err != nil {
		return err
	}

	return nil
}

func (tx Transaction) FindSampleByID(ctx context.Context, id string) (Sample, error) {
	var sample Sample

	if err := tx.dbtx.QueryRow(ctx, `SELECT id, name FROM samples WHERE id = $1`, id).Scan(&sample.ID, &sample.Name); err != nil {
		return Sample{}, err
	}

	return sample, nil
}

func (tx Transaction) ListSamples(ctx context.Context) ([]Sample, error) {
	rows, err := tx.dbtx.Query(ctx, `SELECT id, name FROM samples`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var samples []Sample
	for rows.Next() {
		var sample Sample
		if err := rows.Scan(&sample.ID, &sample.Name); err != nil {
			return nil, err
		}
		samples = append(samples, sample)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return samples, nil
}

func (tx Transaction) DeleteSamples(ctx context.Context) error {
	if _, err := tx.dbtx.Exec(ctx, `DELETE FROM samples`); err != nil {
		return err
	}

	return nil
}
