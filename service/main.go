package main

import (
	"os"
)

func main() {
	kafka := KafkaBroker{
		brokers: os.Getenv("KAFKA_BROKERS"),
	}
	http := HttpHandler{
		kafka: &kafka,
	}
	http.Handle()
}
