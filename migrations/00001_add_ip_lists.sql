-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS whitelist (
    subnet CIDR NOT NULL
);
CREATE TABLE IF NOT EXISTS blacklist (
    subnet CIDR NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS whitelist
DROP TABLE IF EXISTS blacklist
-- +goose StatementEnd
