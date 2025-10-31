
CREATE EXTENSION IF NOT EXISTS timescaledb;

CREATE TABLE IF NOT EXISTS mqtt_data (
    time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    topic TEXT NOT NULL,
    value JSONB,
    pod_name TEXT
);


WITH (
  timescaledb.hypertable,
  timescaledb.partition_column='time',
  timescaledb.segmentby='topic'
);


SELECT create_hypertable('mqtt_data', 'time', if_not_exists => TRUE);

CREATE INDEX IF NOT EXISTS idx_topic ON mqtt_data (topic);