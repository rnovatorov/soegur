BEGIN;

CREATE TABLE es_aggregates (
    id TEXT PRIMARY KEY,
    version INT NOT NULL
);

CREATE TABLE es_events (
    id TEXT PRIMARY KEY,
    sequence_number BIGINT UNIQUE,
    aggregate_id TEXT NOT NULL REFERENCES es_aggregates (id),
    aggregate_version INT NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    metadata JSONB NOT NULL,
    data JSONB NOT NULL,
    UNIQUE (aggregate_id, aggregate_version)
);

CREATE INDEX ON es_events (aggregate_version) INCLUDE (id)
WHERE
    sequence_number IS NULL;

CREATE TABLE es_subscriptions (
    id TEXT PRIMARY KEY,
    position BIGINT NOT NULL DEFAULT 0
);

CREATE TABLE es_subscription_backlogs (
    subscription_id TEXT NOT NULL,
    event_id TEXT NOT NULL REFERENCES es_events (id),
    PRIMARY KEY (subscription_id, event_id)
);

END;
