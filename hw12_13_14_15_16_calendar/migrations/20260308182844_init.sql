-- +goose Up
create table if not exists events (
    id            UUID PRIMARY KEY,
    title         TEXT NOT NULL,
    event_time    TIMESTAMPTZ NOT NULL,
    duration      BIGINT NOT NULL,
    description   TEXT,
    owner_id      UUID NOT NULL,
    notify_before BIGINT
);

-- +goose Down
drop table events;
