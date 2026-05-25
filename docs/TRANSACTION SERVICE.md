Table Transaction 
CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    invoice_number VARCHAR(100),
    user_id UUID,
    amount NUMERIC,
    status VARCHAR(20),
    payment_method VARCHAR(50),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

Endpoint
POST /api/v1/transactions
GET  /api/v1/transactions
GET  /api/v1/transactions/:id