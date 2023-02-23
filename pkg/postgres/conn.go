package postgres

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Conn is a wrapper around pgxpool.Conn
type Conn struct {
	*pgxpool.Conn
}

// Begin wraps conn.Begin
func (c *Conn) Begin(ctx context.Context) (*Tx, error) {
	if err := getChaosError(); err != nil {
		return nil, err
	}
	tx, err := c.Conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &Tx{
		Tx: tx,
	}, nil
}

func (c *Conn) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	if err := getChaosError(); err != nil {
		return nil, err
	}
	return c.Conn.Exec(ctx, sql, args...)
}

func (c *Conn) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	if err := getChaosError(); err != nil {
		return nil, err
	}
	return c.Query(ctx, sql, args...)
}

func (c *Conn) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return c.Conn.QueryRow(ctx, sql, args...)
}

func (c *Conn) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	return c.Conn.SendBatch(ctx, b)
}
