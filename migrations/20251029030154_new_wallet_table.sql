-- +goose Up
-- +goose StatementBegin
CREATE TABLE wallets (
    uuid UUID,
    balance NUMERIC(16, 2) DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS wallets;
-- +goose StatementEnd
