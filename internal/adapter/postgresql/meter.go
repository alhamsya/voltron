package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"

	modelPostgresql "github.com/alhamsya/voltron/internal/core/domain/postgresql"
	modelRequest "github.com/alhamsya/voltron/internal/core/domain/request"
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

func (db *PostgreSQL) GetTimeSeries(ctx context.Context, param *modelRequest.TimeSeries) ([]modelPostgresql.PowerMeter, error) {
	query := `
        SELECT time, value
		FROM power_meter
		WHERE device_id = @device_id
		  AND metric = @metric
		  AND time >= @from
		  AND time < @to
		ORDER BY time ASC
		LIMIT @limit;
    `

	args := pgx.NamedArgs{
		"device_id": param.DeviceID,
		"metric":    param.Metric,
		"from":      param.From,
		"to":        param.To,
		"limit":     5000,
	}

	rows, err := db.Replica.Query(ctx, query, args)
	if err != nil {
		return nil, errors.Wrap(err, "failed Query")
	}
	defer rows.Close()

	var powerMeters []modelPostgresql.PowerMeter
	for rows.Next() {
		var powerMeter modelPostgresql.PowerMeter
		err = rows.Scan(&powerMeter.Time, &powerMeter.Value)
		if err != nil {
			return nil, errors.Wrap(err, "failed rows Scan")
		}
		powerMeters = append(powerMeters, powerMeter)
	}
	return powerMeters, nil
}
