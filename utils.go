package main

import (
	"crypto/sha256"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

func loadRateLimitLuaScript() {
	scriptContent, _ := os.ReadFile(tokenBucketScript)
	sha, _ := rdb.ScriptLoad(ctx, string(scriptContent)).Result()
	logger.Infof("Rate limit Lua script loaded with hash: %s", sha)
	rateLimitSha = sha
}

func initLogger() {
	logger = logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339Nano,
	})
	logger.Info("Logger initialized!")
}

func hashData(data []string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(strings.Join(data, ""))))
}

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	logger.Info("Redis connection successful")
	loadRateLimitLuaScript()
}
