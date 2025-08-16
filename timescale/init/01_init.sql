-- Crea estensione TimescaleDB
CREATE EXTENSION IF NOT EXISTS timescaledb;

-- Tabella semplice per tutti i dati MQTT
CREATE TABLE IF NOT EXISTS mqtt_data (
    time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    topic TEXT NOT NULL,
    value JSONB,
    pod_name TEXT
);

-- Converti in hypertable (il minimo per TimescaleDB)
SELECT create_hypertable('mqtt_data', 'time', if_not_exists => TRUE);

-- Un indice base per i topic
CREATE INDEX IF NOT EXISTS idx_topic ON mqtt_data (topic);