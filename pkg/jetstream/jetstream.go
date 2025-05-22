package jetstream

import (
	"discord/config"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
)

func InitJetstream() (nats.JetStreamContext, error) {
	nc, err := nats.Connect(fmt.Sprintf("nats://%s:%s", config.NatsHost, config.NatsPort))
	if err != nil {
		return nil, err
	}

	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return nil, err
	}

	err = CreateStream(js)
	if err != nil {
		return nil, err
	}

	return js, nil
}

func CreateStream(jetStream nats.JetStreamContext) error {
	_, err := jetStream.StreamInfo(config.JetstreamName)
	if err != nil {
		log.Printf("Stream does not exist, creating stream: %s\n", config.JetstreamName)

		_, err = jetStream.AddStream(
			&nats.StreamConfig{
				Name:       config.JetstreamName,
				Subjects:   []string{"event.>"},
				MaxBytes:   1024 * 1024 * 512,
				Storage:    nats.MemoryStorage,
				Retention:  nats.InterestPolicy,
				Discard:    nats.DiscardOld,
				MaxMsgs:    1000000,
				MaxMsgSize: 512,
				Replicas:   1,
				DenyPurge:  true,
				DenyDelete: true,
			})
		if err != nil {
			return err
		}
	} else {
		log.Printf("Stream already exists: %s\n", config.JetstreamName)
	}
	return nil
}
