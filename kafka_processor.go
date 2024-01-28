package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"github.com/IBM/sarama"
	"io"
	"os"
	"time"
)

type ConsumerGroupHandler struct {
	Pool *WorkerPool
}

func initKafkaConfig() {
	config := sarama.NewConfig()
	config.Version = sarama.V3_6_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Return.Errors = true
	logger.Info("Kafka config loaded")
	kafkaConfig = config
}

func (h *ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	// Setup logic if needed
	return nil
}

func (h *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	// Cleanup logic if needed
	return nil
}

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		filePath := string(msg.Value)
		logger.Infof("Received Kafka message with file path: %s at offset: %d", filePath, msg.Offset)

		if processCSVFile(filePath, h.Pool) {
			logger.Infof("Successfully processed message at offset: %d", msg.Offset)
			session.MarkMessage(msg, "") // Mark message as processed
		} else {
			logger.Errorf("ALERT::An error occurred while processing message with file path: %s at offset: %d",
				filePath, msg.Offset)
		}
	}
	return nil
}

func processCSVFile(filePath string, pool *WorkerPool) bool {
	pool.Reset() // Reset the worker pool for each new message

	file, openErr := os.Open("/Users/anindya.mahajan/Plotline/watch-folder/" + filePath)
	if openErr != nil {
		logger.Errorf("Error opening file: %v\n", openErr)
		return false
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	// Read and discard the first record (header)
	if _, readErr := reader.Read(); readErr != nil {
		logger.Errorf("Error reading header from CSV file: %v", readErr)
		return false
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Errorf("Error reading CSV row: %v\n", err)
			pool.Success = false
			continue
		}

		pool.Work <- row // Send row to the worker pool
	}

	pool.WG.Wait()
	return pool.Success
}

func processRow(data []string) bool {
	hash := hashData(data)
	exists, jsonMarshalErr := rdb.SIsMember(ctx, "processed_hashes", hash).Result()
	if jsonMarshalErr != nil {
		logger.Errorf("Error checking hash existence: %v\n", jsonMarshalErr)
		return false
	}
	if exists {
		logger.Infof("Skipping row with data %v as it has been already processed", data)
		return true
	}

	tokenAcquired := acquireRateLimitToken()
	retry := 1
	for {
		if tokenAcquired {
			break
		} else if retry > maxRetry {
			logger.Infof("Rate limit exceeded, retries exhausted!")
			return false
		} else {
			logger.Infof("Rate limit exceeded, retry number %d", retry)
			time.Sleep(1 * time.Second)
			tokenAcquired = acquireRateLimitToken()
			retry++
		}
	}

	// Create the request body
	reqBody := map[string]string{
		"userID":       data[0],
		"alertMessage": data[1],
	}
	jsonBody, jsonMarshalErr := json.Marshal(reqBody)
	if jsonMarshalErr != nil {
		logger.WithError(jsonMarshalErr).Error("Error marshaling JSON")
		return false
	}

	if makeAPIReqErr := sendPOSTRequest(apiURL, jsonBody); makeAPIReqErr != nil {
		logger.WithError(makeAPIReqErr).Error("Error sending POST request")
		return false
	}

	if redisAddErr := rdb.SAdd(ctx, "processed_hashes", hash).Err(); redisAddErr != nil {
		logger.Errorf("Error adding hash to Redis: %v\n", redisAddErr)
	}
	return true
}
