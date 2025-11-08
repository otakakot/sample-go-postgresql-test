package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	pool *pgxpool.Pool
}

func NewPool(dsn string) (*pgxpool.Pool, error) {
	ctx := context.Background()

	conn, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, conn)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}

func NewDatabase(pool *pgxpool.Pool) (Database, error) {
	return Database{pool: pool}, nil
}

func (db Database) InsertSample(ctx context.Context, name string) (Sample, error) {
	var id string

	if err := db.pool.QueryRow(ctx, `INSERT INTO samples (name) VALUES ($1) RETURNING id`, name).Scan(&id); err != nil {
		return Sample{}, err
	}

	return Sample{ID: id, Name: name}, nil
}

func (db Database) UpdateSample(ctx context.Context, id string, name string) error {
	if _, err := db.pool.Exec(ctx, `UPDATE samples SET name = $1 WHERE id = $2`, name, id); err != nil {
		return err
	}

	return nil
}

func (db Database) FindSampleByID(ctx context.Context, id string) (Sample, error) {
	var sample Sample

	if err := db.pool.QueryRow(ctx, `SELECT id, name FROM samples WHERE id = $1`, id).Scan(&sample.ID, &sample.Name); err != nil {
		return Sample{}, err
	}

	return sample, nil
}

func (db Database) ListSamples(ctx context.Context) ([]Sample, error) {
	rows, err := db.pool.Query(ctx, `SELECT id, name FROM samples`)
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

func (db Database) DeleteSamples(ctx context.Context) error {
	if _, err := db.pool.Exec(ctx, `DELETE FROM samples`); err != nil {
		return err
	}

	return nil
}
