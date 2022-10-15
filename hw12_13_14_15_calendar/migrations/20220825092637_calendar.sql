-- +goose Up
-- +goose StatementBegin
alter table events
    add column notify_after timestamp(0) default NULL::timestamp without time zone,
    add column notified_at timestamp(0) default NULL::timestamp without time zone;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table  events
    drop column notify_after,
    drop column notified_at;
-- +goose StatementEnd
