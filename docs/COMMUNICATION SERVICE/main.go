package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

type PaymentMessage struct {
	Invoice string `json:"invoice"`
	Amount  int    `json:"amount"`
	Phone   string `json:"phone"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	rabbitURL := os.Getenv("RABBITMQ_URL")

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatal("Failed connect RabbitMQ:", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	queueName := "payment.paid"

	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Waiting message...")

	forever := make(chan bool)

	go func() {
		for d := range msgs {

			var msg PaymentMessage

			err := json.Unmarshal(d.Body, &msg)
			if err != nil {
				log.Println("Invalid JSON:", err)
				continue
			}

			log.Printf("Receive payment: %+v\n", msg)

			err = sendWhatsApp(msg)
			if err != nil {
				log.Println("Send WA failed:", err)
				continue
			}

			log.Println("WhatsApp sent")
		}
	}()

	<-forever
}

func sendWhatsApp(msg PaymentMessage) error {

	message := fmt.Sprintf(
		"Pembayaran berhasil diterima.\n\nInvoice: %s\nNominal: Rp%d\n\nTerima kasih.",
		msg.Invoice,
		msg.Amount,
	)

	payload := map[string]string{
		"target":  msg.Phone,
		"message": message,
	}

	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest(
		"POST",
		os.Getenv("WA_API_URL"),
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", os.Getenv("WA_API_TOKEN"))

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	log.Println("WA Status:", resp.Status)

	return nil
}