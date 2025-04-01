package processor

import (
	"discord/config"
	"discord/internal/methods"
	"discord/models"
	"fmt"
	"github.com/nats-io/nats.go"
)

func (h *Handler) MessageProcessor() {
	if h.js == nil {
		fmt.Println("Jetstream is nil")
		return
	}
	fmt.Println("Message Processor started")
	_, err := h.js.Jetstream.Subscribe(
		config.JetstreamSubjectsEventMsg, func(msg *nats.Msg) {
			data, err := methods.StringToStruct[models.Message](string(msg.Data))
			if err != nil {
				fmt.Println(err)
			}
			err = h.ch.BatchMessage(&data)
			if err != nil {
				fmt.Println(err)
			}
			err = msg.Ack()
			if err != nil {
				return
			}
		}, nats.ManualAck(), nats.Durable("backend-message-consumer"))
	if err != nil {
		fmt.Println("Error subscribing to Jetstream:", err)
	}
}
