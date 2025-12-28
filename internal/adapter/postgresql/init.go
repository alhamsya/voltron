package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

const DriverPostgres = "postgres"

type Config struct {
	Username        string
	Password        string
	Host            string
	Port            int
	Name            string
	MaxConns        int
	MinIdleConns    int
	MaxConnLifetime int
	MaxConnIdleTime int
}

type PostgreSQL struct {
	Primary *pgxpool.Pool
	Replica *pgxpool.Pool
}

func New(primary, replica *pgxpool.Pool) *PostgreSQL {
	if primary == nil || replica == nil {
		panic(errors.New("primary or replica param is nil"))
	}
	return &PostgreSQL{
		Primary: primary,
		Replica: replica,
	}
}

func Connect(ctx context.Context, cfg *Config) *pgxpool.Pool {
	connString := fmt.Sprintf("%s://%s:%s@%s:%d/%s",
		DriverPostgres,
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		panic(errors.Wrap(err, "failed pgxpool ParseConfig"))
	}

	db, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		panic(errors.Wrap(err, "failed pgxpool NewWithConfig"))
	}

	return db
}

func (db *PostgreSQL) Ping(ctx context.Context) error {
	return db.Ping(ctx)
}

func (db *PostgreSQL) Close() {
	db.Close()
}
