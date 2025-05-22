package worker

import (
	"discord/config"
	"discord/internal/methods"
	"discord/models"
	"fmt"
	"github.com/nats-io/nats.go"
)

func (w *Worker) MessageProcessor() {

	if w.stream.Jetstream == nil {
		fmt.Println("Jetstream is nil")
		return
	}
	fmt.Println("Message Processor started")
	_, err := w.stream.Jetstream.Subscribe(
		config.JetstreamSubjectEventMsg, func(msg *nats.Msg) {
			data, err := methods.StringToStruct[models.Message](string(msg.Data))
			if err != nil {
				fmt.Println(err)
			}
			err = w.batcher.BatchMessage(&data)
			if err != nil {
				fmt.Println(err)
			}
			err = msg.Ack()
			if err != nil {
				fmt.Println("Error acknowledging message:", err)
			}
		}, nats.ManualAck(), nats.Durable("backend-message-consumer"))
	if err != nil {
		fmt.Println("Error subscribing to Jetstream:", err)
	}
}
