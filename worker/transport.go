package main

import (
	"context"
	"encoding/json"

	"github.com/Shopify/sarama"
	"github.com/ricardo-ch/go-kafka"
)

func makePostCreatedEndpoint(s ClasspageService) kafka.Handler {
	return func(ctx context.Context, msg *sarama.ConsumerMessage) error {
		postCreatedMsg := PostCreatedMsg{}
		err := json.Unmarshal(msg.Value, &postCreatedMsg)
		if err != nil {
			return err
		}
		return s.OnPostCreated(postCreatedMsg)
	}
}
