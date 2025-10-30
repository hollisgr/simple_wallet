-- +goose Up
-- +goose StatementBegin
CREATE TABLE wallets (
    uuid UUID,
    balance NUMERIC(16, 2) DEFAULT 0.00
);

CREATE INDEX wallet_uuid_idx ON public.wallets(uuid);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS wallets;
-- +goose StatementEnd
