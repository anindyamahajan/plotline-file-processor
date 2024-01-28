package main

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	kafkaBroker          = "localhost:9092"
	kafkaTopic           = "csv_file_address"
	maxRetry             = 4
	initialBackoff       = 1000 * time.Millisecond
	backoffFactor        = 3
	redisAddr            = "localhost:6379"
	apiURL               = "http://localhost:9900/webhook/"
	apiRequestTimeout    = 10 * time.Second
	tokenBucketScript    = "token_script.lua"
	tokenBucketKey       = "api_rate_limit"
	tokenBucketRate      = 10 // Number of tokens added per interval
	tokenBucketInterval  = 1 * time.Second
	tokenBucketMaxTokens = 10
	workers              = 10
)

var (
	rdb          *redis.Client
	rateLimitSha string
	logger       *logrus.Logger
	ctx          context.Context
	kafkaConfig  *sarama.Config
)
