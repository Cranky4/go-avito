-- +goose Up
-- +goose StatementBegin
create table events
(
    id                   uuid         not null
        primary key,
    title               varchar(255)  not null,
    starts_at           timestamp(0)  default NULL::timestamp without time zone,
    ends_at             timestamp(0)  default NULL::timestamp without time zone
);

create index idx_events_starts_at
    on events (starts_at);
create index idx_events_ends_at
    on events (ends_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd
