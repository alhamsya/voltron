package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"

	modelPostgresql "github.com/alhamsya/voltron/internal/core/domain/postgresql"
)

func (db *PostgreSQL) CreatePowerMeterReading(ctx context.Context, pmReading *modelPostgresql.PowerMeter) error {
	query := `
        INSERT INTO books (title, author, quantity)
		VALUES (@time, @device_id, @metric, @value, @event_Hash)
    `

	// Define the named arguments for the query.
	args := pgx.NamedArgs{
		"time":       pmReading.Time,
		"device_id":  pmReading.DeviceID,
		"metric":     pmReading.Metric,
		"value":      pmReading.Value,
		"event_Hash": pmReading.EventHash,
	}

	_, err := db.Primary.Exec(ctx, query, args)
	if err != nil {
		return errors.Wrap(err, "failed exec")
	}

	return nil
}

func (db *PostgreSQL) BulkPowerMeter(ctx context.Context, param []modelPostgresql.PowerMeter) error {
	return pgx.BeginFunc(ctx, db.Primary, func(tx pgx.Tx) error {
		const q = `
			INSERT INTO power_meter (time, device_id, metric, value, event_hash)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT DO NOTHING;
		`

		b := &pgx.Batch{}
		for _, r := range param {
			b.Queue(q, r.Time, r.DeviceID, r.Metric, r.Value, r.EventHash)
		}

		br := tx.SendBatch(ctx, b)
		defer br.Close()

		for i := 0; i < len(param); i++ {
			if _, err := br.Exec(); err != nil {
				return fmt.Errorf("batch exec #%d: %w", i, err)
			}
		}
		return nil
	})
}
