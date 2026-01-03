package postgresql

import (
	"context"
	"fmt"
	"github.com/alhamsya/voltron/internal/core/domain/constant"
	modelPower "github.com/alhamsya/voltron/internal/core/domain/power"
	"time"

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

func (db *PostgreSQL) GetMeterTimeSeries(ctx context.Context, param *modelRequest.TimeSeries) ([]modelPostgresql.PowerMeter, error) {
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

func (db *PostgreSQL) GetMeterLatestMeter(ctx context.Context, deviceID string) ([]modelPostgresql.Latest, error) {
	const query = `
			SELECT DISTINCT ON (metric)
			  metric,
			  time,
			  value
			FROM power_meter
			WHERE device_id = $1
			ORDER BY metric, time DESC
		`

	rows, err := db.Replica.Query(ctx, query, deviceID)
	if err != nil {
		return nil, errors.Wrap(err, "failed Query")
	}
	defer rows.Close()

	var resp []modelPostgresql.Latest
	for rows.Next() {
		var l modelPostgresql.Latest
		if err = rows.Scan(&l.Metric, &l.Time, &l.Value); err != nil {
			return nil, errors.Wrap(err, "failed rows scan")
		}
		resp = append(resp, l)
	}

	return resp, nil
}

func (db *PostgreSQL) GetMeterDailyUsage(ctx context.Context, deviceID string, from, to time.Time) ([]modelPostgresql.DailyUsage, error) {
	const query = `
			SELECT
			  day,
			  usage_kwh
			FROM daily_usage
			WHERE device_id = $1
			  AND day >= $2::date
			  AND day <  $3::date
			ORDER BY day ASC;
		`

	rows, err := db.Replica.Query(ctx, query, deviceID, from, to)
	if err != nil {
		return nil, errors.Wrap(err, "failed Query")
	}
	defer rows.Close()

	var out []modelPostgresql.DailyUsage
	for rows.Next() {
		var day time.Time
		var usage float64
		if err = rows.Scan(&day, &usage); err != nil {
			return nil, errors.Wrap(err, "failed rows scan")
		}
		out = append(out, modelPostgresql.DailyUsage{
			Day:      day.Format(constant.DateOnly),
			UsageKwh: usage,
		})
	}

	return out, nil
}

func (db *PostgreSQL) GetBillingSummary(ctx context.Context, deviceID string, from, to time.Time) (float64, error) {
	const q = `
		SELECT COALESCE(SUM(usage_kwh), 0) AS total_kwh
		FROM daily_usage
		WHERE device_id = $1
		  AND day >= $2::date
		  AND day <  $3::date;
		`
	var total float64
	err := db.Replica.QueryRow(ctx, q, deviceID, from, to).Scan(&total)
	return total, err
}

func (db *PostgreSQL) GetDailyUsageLines(ctx context.Context, deviceID string, from, to time.Time) ([]modelPower.DailyLine, error) {
	const q = `
		SELECT day, usage_kwh
		FROM daily_usage
		WHERE device_id = $1
		  AND day >= $2::date
		  AND day <  $3::date
		ORDER BY day ASC;
		`
	rows, err := db.Replica.Query(ctx, q, deviceID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []modelPower.DailyLine
	for rows.Next() {
		var day time.Time
		var u float64
		if err := rows.Scan(&day, &u); err != nil {
			return nil, err
		}
		out = append(out, modelPower.DailyLine{Day: day.Format(constant.DateOnly), UsageKwh: u})
	}
	return out, rows.Err()
}
