-- execute 1
CREATE EXTENSION IF NOT EXISTS timescaledb;

DROP TABLE IF EXISTS power_meter;

CREATE TABLE power_meter (
  time       TIMESTAMPTZ NOT NULL,
  device_id  TEXT NOT NULL,
  metric     TEXT NOT NULL,
  value      NUMERIC(18,6) NOT NULL,
  event_hash TEXT NOT NULL,
  PRIMARY KEY (device_id, time, metric, event_hash)
);

SELECT create_hypertable('power_meter', 'time', if_not_exists => TRUE);

CREATE INDEX idx_power_meter_device_time
  ON power_meter (device_id, time DESC);

CREATE INDEX idx_power_meter_metric_time
  ON power_meter (metric, time DESC);

-- execute 2

DROP MATERIALIZED VIEW IF EXISTS daily_usage;

CREATE MATERIALIZED VIEW daily_usage
WITH (timescaledb.continuous) AS
SELECT
  time_bucket('1 day', time) AS day,
  device_id,
  GREATEST(last(value, time) - first(value, time), 0) AS usage_kwh
FROM power_meter
WHERE metric = 'total_import_kwh'
GROUP BY day, device_id;

-- execute 3
-- (optional) kalau sebelumnya sudah pernah dibuat, drop dulu biar gak error
SELECT remove_continuous_aggregate_policy('daily_usage')
WHERE EXISTS (
  SELECT 1
  FROM timescaledb_information.jobs j
  JOIN timescaledb_information.job_stats js ON js.job_id = j.job_id
  WHERE j.proc_name = 'policy_refresh_continuous_aggregate'
);

SELECT add_continuous_aggregate_policy('daily_usage',
  start_offset => INTERVAL '7 days',
  end_offset   => INTERVAL '1 day',
  schedule_interval => INTERVAL '30 minutes'
);

CALL refresh_continuous_aggregate(
  'daily_usage',
  '2025-12-19 00:00:00+00',
  '2025-12-20 00:00:00+00'
);