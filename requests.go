package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

func acquireRateLimitToken() bool {
	result, err := rdb.EvalSha(ctx, rateLimitSha, []string{tokenBucketKey}, tokenBucketRate, tokenBucketInterval.Milliseconds(), tokenBucketMaxTokens).Result()
	if err != nil {
		logger.Errorf("Error executing Lua script: %v\n", err)
		return false
	}

	return result.(int64) > 0
}

func sendPOSTRequest(url string, data []byte) error {
	retryCount := 0
	backoff := initialBackoff

	client := &http.Client{
		Timeout: apiRequestTimeout,
	}

	for {
		resp, err := client.Post(url, "application/json", bytes.NewBuffer(data))
		if err != nil {
			if retryCount >= maxRetry {
				return fmt.Errorf("failed to send POST request after %d retries: %v", maxRetry, err)
			}
			logger.Errorf("Could not make API call due to error: %v. Retrying.", err)
		} else {
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				if retryCount >= maxRetry {
					return fmt.Errorf("failed to send POST request after %d retries: received status code %d", maxRetry, resp.StatusCode)
				}
				logger.Errorf("Could not make API call. Received status code %d. Retrying.", resp.StatusCode)
			} else {
				logger.Info("Successfully posted API request")
				return nil
			}
		}

		time.Sleep(backoff)
		retryCount++
		backoff *= time.Duration(backoffFactor)
	}
}
