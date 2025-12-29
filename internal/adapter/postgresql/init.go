package postgresql

import (
	"context"
	"fmt"
	"github.com/alhamsya/voltron/pkg/manager/config"

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
	Cfg *config.Application

	Primary *pgxpool.Pool
	Replica *pgxpool.Pool
}

func New(cfg *config.Application, primary, replica *pgxpool.Pool) *PostgreSQL {
	if primary == nil || replica == nil {
		panic(errors.New("primary or replica param is nil"))
	}
	return &PostgreSQL{
		Cfg:     cfg,
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

	if err = db.Ping(ctx); err != nil {
		panic(errors.Wrap(err, "failed ping database"))
	}

	return db
}

func (db *PostgreSQL) Ping(ctx context.Context) error {
	return db.Ping(ctx)
}

func (db *PostgreSQL) Close() {
	db.Close()
}
