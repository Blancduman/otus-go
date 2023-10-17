-- +goose Up
CREATE TABLE event (
    id serial primary key,
    title text,
    start_date timestamptz not null,
    end_date timestamptz not null,
    description text,
    owner_id bigint,
    remind_in timestamptz
);

CREATE INDEX start_datex ON event (start_date);

-- +goose Down
DROP TABLE event;