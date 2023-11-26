CREATE TABLE IF NOT EXISTS event (
                       id serial primary key,
                       title text,
                       start_date timestamp with time zone not null,
                       end_date timestamp with time zone not null,
                       description text,
                       owner_id bigint,
                       remind_in timestamp with time zone
);

CREATE INDEX start_datex ON event (start_date);