CREATE TABLE transactions (
    transaction_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    account_id BIGINT NOT NULL,
    operation_type_id BIGINT NOT NULL,
    amount NUMERIC(12,2) NOT NULL,
    event_date TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_account FOREIGN KEY(account_id) REFERENCES accounts(account_id),
    CONSTRAINT fk_operation_type FOREIGN KEY(operation_type_id) REFERENCES operation_types(operation_type_id)
);