project/
├── cmd/
│   └── app/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── delivery/
│   │   └── http/
│   │       ├── handler/
│   │       │   ├── auth_handler.go
│   │       │   └── transaction_handler.go
│   │       └── routes.go
│   ├── entity/
│   │   ├── user.go
│   │   └── transaction.go
│   ├── helper/
│   │   ├── jwt.go
│   │   └── rabbitmq.go
│   ├── middleware/
│   │   └── jwt.go
│   ├── repository/
│   │   ├── user_repository.go
│   │   └── transaction_repository.go
│   ├── usecase/
│   │   ├── auth_usecase.go
│   │   └── transaction_usecase.go
│   └── database/
│       └── postgres.go
├── go.mod
└── .env    