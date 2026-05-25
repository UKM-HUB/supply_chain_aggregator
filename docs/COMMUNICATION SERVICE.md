RabbitMQ Consumer
RabbitMQ Queue:
payment.paid

Flow

Transaction PAID
   ↓
RabbitMQ Publish
   ↓
Notification Service Consume
   ↓
Send WhatsApp

Contoh Message
{
  "invoice": "INV-001",
  "amount": 1000000,
  "phone": "628123456789"
}

WhatsApp Message

Pembayaran berhasil diterima.

Invoice: INV-001
Nominal: Rp1.000.000

Terima kasih.