package main

import (
    "encoding/json"
    "net/http"
    "strings"
    "time"

    "github.com/google/uuid"
)

func CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
    var t Transaction

    err := json.NewDecoder(r.Body).Decode(&t)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    now := time.Now()

    t.ID = uuid.New().String()
    t.CreatedAt = now
    t.UpdatedAt = now

    err = CreateTransaction(t)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(t)
}

func GetTransactionsHandler(w http.ResponseWriter, r *http.Request) {
    data, err := GetTransactions()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(data)
}

func GetTransactionByIDHandler(w http.ResponseWriter, r *http.Request) {
    id := strings.TrimPrefix(r.URL.Path, "/api/v1/transactions/")

    data, err := GetTransactionByID(id)
    if err != nil {
        http.Error(w, "not found", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(data)
}