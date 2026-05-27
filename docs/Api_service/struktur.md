Struktur Project
project/
├── cmd/
│   └── main.go
├── config/
│   └── database.go
├── routes/
│   └── routes.go
├── middleware/
│   └── jwt.go
├── models/
│   ├── user.go
│   ├── transaction.go
│   └── report.go
├── controllers/
│   ├── auth_controller.go
│   ├── user_controller.go
│   ├── transaction_controller.go
│   ├── report_controller.go
│   ├── nearby_controller.go
│   ├── communication_controller.go
│   └── gateway_controller.go
├── services/
│   ├── auth_service.go
│   ├── user_service.go
│   ├── transaction_service.go
│   └── report_service.go
├── repositories/
│   ├── user_repository.go
│   └── transaction_repository.go
├── utils/
│   ├── response.go
│   └── jwt.go
├── .env
├── go.mod
└── go.sum