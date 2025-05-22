package worker

import (
	"discord/config"
	"discord/internal/methods"
	"discord/models"
	redisdb "discord/pkg/redis"
	"fmt"
	"github.com/nats-io/nats.go"
)

func (w *Worker) UserProcessor() {
	if w.stream.Jetstream == nil {
		return
	}
	fmt.Println("User Processor started")
	_, err := w.stream.Jetstream.Subscribe(
		config.JetstreamSubjectUserProfile, func(msg *nats.Msg) {
			data, err := methods.StringToStruct[models.UserProfile](string(msg.Data))
			if err != nil {
				fmt.Println(err)
			}
			redisdb.HashSet(1, "user-profile-hash", data.UserId, data.AvatarHash)
			err = msg.Ack()
			if err != nil {
				fmt.Println("Error acknowledging message:", err)
			}
		}, nats.ManualAck(), nats.Durable("backend-user-pfp-consumer"))
	if err != nil {
		fmt.Println("Error subscribing to Jetstream:", err)
	}
}
