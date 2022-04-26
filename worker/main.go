package main

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ricardo-ch/go-kafka"
)

var (
	serviceName   = "vn.freeclass.classpage-worker"
	vendors       = VendorManager{}
	kafkaBrokers  []string
	auroraConfigs map[string]string
)

func init() {
	godotenv.Load()
	kafkaBrokers = strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	auroraConfigs = map[string]string{
		"db_user": os.Getenv("DB_USER"),
		"db_pass": os.Getenv("DB_PASS"),
		"db_host": os.Getenv("DB_HOST"),
		"db_name": os.Getenv("DB_NAME"),
	}
}

func main() {
	var (
		logger        = log.NewLogfmtLogger(os.Stderr)
		aurora        = vendors.NewAuroraClient(auroraConfigs)
		awsChime      = vendors.NewChimeClient()
		kafkaProducer = vendors.NewKafkaProducer(kafkaBrokers)
	)
	classpageRepository := NewClasspageRepository(aurora, awsChime, kafkaProducer)
	classpageService := NewClasspageService(classpageRepository)

	// HTTP server
	go func() {
		r := mux.NewRouter()
		r.Handle("/metrics", promhttp.Handler())
		logger.Log("msg", "HTTP server is listening", "addr", ":8080")
		logger.Log("err", http.ListenAndServe(":8080", r))
	}()

	// Kafka consumer
	handlers := kafka.Handlers{}
	handlers["classpages.fct.post-created.0"] = makePostCreatedEndpoint(classpageService)

	listener, err := kafka.NewListener(kafkaBrokers, serviceName, handlers, kafka.WithInstrumenting())
	if err != nil {
		logger.Log("msg", "Could not initialise listener", "err", err)
		return
	}
	defer listener.Close()
	err = listener.Listen(context.Background())
	if err != nil {
		logger.Log("msg", "Listener closed with error", "err", err)
	}
	logger.Log("msg", "Listener stopped")
}
