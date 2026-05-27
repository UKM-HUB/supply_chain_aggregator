package main

import (
    "log"
    "net/http"
)

func main() {
    InitDB()

    http.HandleFunc("/api/v1/transactions", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            CreateTransactionHandler(w, r)
        } else if r.Method == http.MethodGet {
            GetTransactionsHandler(w, r)
        }
    })

    http.HandleFunc("/api/v1/transactions/", GetTransactionByIDHandler)

    log.Println("Server running on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}