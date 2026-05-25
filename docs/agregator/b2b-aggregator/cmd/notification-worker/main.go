package main

import (
	"log"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Gagal connect ke RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Gagal buka channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("order_events", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Gagal declare queue: %v", err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Gagal register consumer: %v", err)
	}

	log.Println("Worker Notification berjalan. Menunggu event dari RabbitMQ...")

	// Mencegah program exit
	var forever chan struct{}
	go func() {
		for d := range msgs {
			log.Printf("📥 Menerima Notifikasi Pesanan Baru: %s", d.Body)
			// TODO: Panggil API WhatsApp / Email (misal: SMTP / Twilio) di sini
		}
	}()
	<-forever
}
