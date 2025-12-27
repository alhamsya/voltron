-- execute 1
CREATE EXTENSION IF NOT EXISTS timescaledb;

DROP TABLE IF EXISTS pm_readings;

CREATE TABLE pm_readings (
  time       TIMESTAMPTZ NOT NULL,
  device_id  TEXT NOT NULL,
  metric     TEXT NOT NULL CHECK (metric IN ('Volts','Current','Active_Power','Total_Import_kWh')),
  value      NUMERIC(18,6) NOT NULL,
  event_hash TEXT NOT NULL,
  PRIMARY KEY (device_id, time, metric, event_hash)
);

SELECT create_hypertable('pm_readings', 'time', if_not_exists => TRUE);

CREATE INDEX idx_pm_readings_device_time ON pm_readings (device_id, time DESC);
CREATE INDEX idx_pm_readings_metric_time ON pm_readings (metric, time DESC);

-- execute 2

DROP MATERIALIZED VIEW IF EXISTS daily_usage;

CREATE MATERIALIZED VIEW daily_usage
WITH (timescaledb.continuous) AS
SELECT
  time_bucket('1 day', time) AS day,
  device_id,
  (max(value) - min(value)) AS usage_kwh
FROM pm_readings
WHERE metric = 'Total_Import_kWh'
GROUP BY day, device_id;


-- execute 3
SELECT add_continuous_aggregate_policy('daily_usage',
  start_offset => INTERVAL '7 days',
  end_offset   => INTERVAL '1 day',
  schedule_interval => INTERVAL '30 minutes'
);