CREATE TABLE payment_logs (
    id                UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_number    VARCHAR(100)  NOT NULL REFERENCES transactions (invoice_number) ON DELETE RESTRICT,
    amount            NUMERIC(15,2) NOT NULL,
    user_phone        VARCHAR(20),
    payment_url       TEXT,
    xendit_invoice_id VARCHAR(255),
    status            VARCHAR(20)   NOT NULL DEFAULT 'pending',
    created_at        TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    CONSTRAINT payment_logs_status_check CHECK (status IN ('pending', 'paid', 'failed'))
);

CREATE INDEX idx_payment_logs_invoice_number ON payment_logs (invoice_number);
CREATE INDEX idx_payment_logs_status         ON payment_logs (status);
