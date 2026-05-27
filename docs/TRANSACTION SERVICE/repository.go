package main

import (
    "database/sql"
)

func CreateTransaction(tx Transaction) error {
    query := `
    INSERT INTO transactions (
        id, invoice_number, user_id, amount, status, payment_method, created_at, updated_at
    ) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
    `

    _, err := DB.Exec(query,
        tx.ID,
        tx.InvoiceNumber,
        tx.UserID,
        tx.Amount,
        tx.Status,
        tx.PaymentMethod,
        tx.CreatedAt,
        tx.UpdatedAt,
    )

    return err
}

func GetTransactions() ([]Transaction, error) {
    rows, err := DB.Query(`SELECT id, invoice_number, user_id, amount, status, payment_method, created_at, updated_at FROM transactions`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var results []Transaction

    for rows.Next() {
        var t Transaction
        err := rows.Scan(
            &t.ID,
            &t.InvoiceNumber,
            &t.UserID,
            &t.Amount,
            &t.Status,
            &t.PaymentMethod,
            &t.CreatedAt,
            &t.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }

        results = append(results, t)
    }

    return results, nil
}

func GetTransactionByID(id string) (Transaction, error) {
    var t Transaction

    err := DB.QueryRow(`
        SELECT id, invoice_number, user_id, amount, status, payment_method, created_at, updated_at
        FROM transactions WHERE id=$1
    `, id).Scan(
        &t.ID,
        &t.InvoiceNumber,
        &t.UserID,
        &t.Amount,
        &t.Status,
        &t.PaymentMethod,
        &t.CreatedAt,
        &t.UpdatedAt,
    )

    if err != nil && err != sql.ErrNoRows {
        return t, err
    }

    return t, err
}