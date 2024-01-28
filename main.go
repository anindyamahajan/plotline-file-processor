package main

import (
	"context"
	"github.com/IBM/sarama"
)

func main() {
	ctx = context.Background()
	initLogger()
	initRedis()
	initKafkaConfig()

	consumerGroup, newConsumerErr := sarama.NewConsumerGroup([]string{kafkaBroker}, "csv-processing-group", kafkaConfig)
	if newConsumerErr != nil {
		logger.WithError(newConsumerErr).Fatal("Error creating consumer group")
	}
	defer consumerGroup.Close()

	// Create a worker pool
	pool := NewWorkerPool(workers)
	logger.Info("Worker pool created with workers: ", workers)
	defer close(pool.Work)

	// Start the ConsumerGroupHandler
	handler := &ConsumerGroupHandler{Pool: pool}
	for {
		if consumeErr := consumerGroup.Consume(ctx, []string{kafkaTopic}, handler); consumeErr != nil {
			logger.WithError(consumeErr).Error("Error in consumer group consume")
		}
		if ctx.Err() != nil {
			return
		}
	}
}
