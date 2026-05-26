CREATE TABLE transactions (
    id             UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_number VARCHAR(100)  NOT NULL,
    user_id        UUID          NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    amount         NUMERIC(15,2) NOT NULL,
    status         VARCHAR(20)   NOT NULL DEFAULT 'pending',
    payment_method VARCHAR(50),
    created_at     TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    CONSTRAINT transactions_invoice_unique UNIQUE (invoice_number),
    CONSTRAINT transactions_status_check   CHECK  (status IN ('pending', 'paid', 'failed', 'cancelled')),
    CONSTRAINT transactions_amount_check   CHECK  (amount > 0)
);

CREATE INDEX idx_transactions_user_id        ON transactions (user_id);
CREATE INDEX idx_transactions_status         ON transactions (status);
CREATE INDEX idx_transactions_created_at     ON transactions (created_at);
CREATE INDEX idx_transactions_invoice_number ON transactions (invoice_number);
