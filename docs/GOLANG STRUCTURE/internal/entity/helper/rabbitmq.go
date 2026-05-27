package helper

import (
    "encoding/json"
    "log"

    amqp "github.com/rabbitmq/amqp091-go"
)

var Conn *amqp.Connection

func InitRabbitMQ(url string) {
    conn, err := amqp.Dial(url)
    if err != nil {
        log.Fatal(err)
    }

    Conn = conn
}

func Publish(queue string, body interface{}) error {
    ch, err := Conn.Channel()
    if err != nil {
        return err
    }

    defer ch.Close()

    _, err = ch.QueueDeclare(
        queue,
        true,
        false,
        false,
        false,
        nil,
    )

    if err != nil {
        return err
    }

    data, err := json.Marshal(body)
    if err != nil {
        return err
    }

    return ch.Publish(
        "",
        queue,
        false,
        false,
        amqp.Publishing{
            ContentType: "application/json",
            Body:        data,
        },
    )
}