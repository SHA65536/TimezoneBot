-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS timezones (
    user_id VARCHAR(20) PRIMARY KEY,
    timezone VARCHAR(32) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS timezones;
-- +goose StatementEnd
