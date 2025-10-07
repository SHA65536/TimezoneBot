-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS streams (
    guild_id VARCHAR(20) PRIMARY KEY,
    alert_channel VARCHAR(20) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
