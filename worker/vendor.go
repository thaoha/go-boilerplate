package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/chime"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type VendorManager struct{}

func (v *VendorManager) NewKafkaProducer(brokers []string) sarama.SyncProducer {
	var config = sarama.NewConfig()
	config.Producer.Timeout = 5 * time.Second
	config.Producer.Retry.Max = 3
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Version = sarama.V1_1_1_0
	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		panic(err)
	}
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	return producer
}

func (v *VendorManager) NewAuroraClient(configs map[string]string) *gorm.DB {
	defaultConfigs := map[string]string{
		"db_user": "root",
		"db_pass": "root",
		"db_host": "localhost",
		"db_name": "example",
	}
	for k, v := range configs {
		defaultConfigs[k] = v
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True",
		defaultConfigs["db_user"],
		defaultConfigs["db_pass"],
		defaultConfigs["db_host"],
		defaultConfigs["db_name"],
	)
	client, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return client
}

func (v *VendorManager) NewChimeClient() *chime.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}
	return chime.NewFromConfig(cfg)
}
