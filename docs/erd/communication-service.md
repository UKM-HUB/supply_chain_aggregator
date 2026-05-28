# ERD — Communication Service

```mermaid
erDiagram
    RABBITMQ_QUEUE {
        string queue_name   "payment.paid"
        string routing_key  "payment.paid"
    }

    PAYMENT_PAID_EVENT {
        string invoice "invoice_number from transaction"
        float  amount  "payment amount in IDR"
        string phone   "recipient phone number (E.164)"
    }

    WHATSAPP_MESSAGE {
        string to      "phone number"
        string body    "formatted notification text"
    }

    %% A queue holds one or many unconsumed events at any moment
    RABBITMQ_QUEUE  ||--|{  PAYMENT_PAID_EVENT  : "delivers one or many events"

    %% Each consumed event triggers exactly one WhatsApp message; each message originates from exactly one event
    PAYMENT_PAID_EVENT ||--||  WHATSAPP_MESSAGE : "triggers exactly one notification"
```

## Cardinality rationale
| Relationship | Left | Right | Reason |
|---|---|---|---|
| RABBITMQ_QUEUE → PAYMENT_PAID_EVENT | exactly one | one or many | The queue always has at least one message in flight when the worker is processing; it accumulates many events over time |
| PAYMENT_PAID_EVENT → WHATSAPP_MESSAGE | exactly one | exactly one | Every consumed event results in exactly one WhatsApp notification call |

## Notes
- This service owns **no database table**; it is a pure event consumer.
- Subscribes to the `payment.paid` queue on RabbitMQ (AMQP) with manual ack/nack.
- WhatsApp message format:
  ```
  Pembayaran berhasil diterima.
  Invoice: INV-20240526-001
  Nominal: Rp1.000.000
  Terima kasih.
  ```
- When `RABBITMQ_URL` is empty the worker logs a warning and blocks (graceful fallback for local dev).
- Graceful shutdown via `signal.NotifyContext` (SIGINT / SIGTERM).
