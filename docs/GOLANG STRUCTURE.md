project/
├── cmd/
├── internal/
│   ├── delivery/
│   │   ├── http/
│   │   └── grpc/
│   ├── usecase/
│   ├── repository/
│   ├── entity/
│   ├── middleware/
│   ├── config/
│   └── helper/
├── pkg/
├── migrations/
├── docs/
├── scripts/
└── main.go

Entity Example
User
type User struct {
    ID       uuid.UUID
    Name     string
    Email    string
    Password string
    Role     string
}

Transaction
type Transaction struct {
    ID            uuid.UUID
    InvoiceNumber string
    Amount        float64
    Status        string
    UserID        uuid.UUID
}

RabbitMQ Publisher
func Publish(queue string, body interface{}) error {
    ch, _ := conn.Channel()

    data, _ := json.Marshal(body)

    return ch.Publish(
        "",
        queue,
        false,
        false,
        amqp.Publishing{
            ContentType: "application/json",
            Body: data,
        },
    )
}

Middleware JWT
e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
    SigningKey: []byte("secret"),
}))

Arsitektur Final

                +-------------------+
                |     Frontend      |
                +-------------------+
                          |
                          v
                +-------------------+
                |    API Gateway    |
                +-------------------+
                          |
      ------------------------------------------------
      |         |           |         |              |
      v         v           v         v              v
   Auth     Transaction   Report   Nearby   Communication
 Service      Service    Service   Service      Service
      |
      v
+-------------------+
| Payment Gateway   |
|      Xendit       |
+-------------------+
      |
      v
+-------------------+
| PostgreSQL        |
+-------------------+

      |
      v
+-------------------+
| RabbitMQ          |
+-------------------+
      |
      v
+-------------------+
| WhatsApp Service  |
+-------------------+

Flow Real Production

Create Order
   ↓
Generate VA
   ↓
Waiting Payment
   ↓
Webhook Paid
   ↓
Update DB
   ↓
Publish Event
   ↓
Send WA
   ↓
Generate Report
   ↓
Dashboard Update