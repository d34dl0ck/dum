package adapters

import (
	"dum/internal/machines/entities"
	"dum/pkg/machines/notification"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

type RabbitNotifictionStrategy struct {
	RabbitUrl string
}

func (s RabbitNotifictionStrategy) Notify(id entities.MachineId, level entities.HealthLevel) error {
	conn, err := amqp.Dial(s.RabbitUrl)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"machine_health", // name
		false,            // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		return err
	}

	msg := notification.HealthChangedEvent{
		Level:     int(level),
		MachineId: uuid.UUID(id),
	}

	raw, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        raw,
		})
	if err != nil {
		return err
	}

	return nil
}
