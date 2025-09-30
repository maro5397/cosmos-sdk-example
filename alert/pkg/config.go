package pkg

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Nodes            []Node
	WebhookURL       string
	PollInterval     time.Duration
	RPCTimeout       time.Duration
	StopDetectWindow time.Duration
}

type Node struct {
	Name             string
	RPC              string
	ValidatorAddress string
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getMs(key string, defaultValue int) time.Duration {
	if value := os.Getenv(key); value != "" {
		if num, _ := strconv.Atoi(value); num > 0 {
			return time.Duration(num) * time.Millisecond
		}
	}
	return time.Duration(defaultValue) * time.Millisecond
}

func getSec(key string, defaultValue int) time.Duration {
	if value := os.Getenv(key); value != "" {
		if num, _ := strconv.Atoi(value); num > 0 {
			return time.Duration(num) * time.Second
		}
	}
	return time.Duration(defaultValue) * time.Second
}

func Load() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return &Config{
		Nodes: []Node{
			{
				Name:             getEnv("REGION", ""),
				RPC:              getEnv("REGION_RPC", ""),
				ValidatorAddress: getEnv("REGION_VALIDATOR_ADDRESS=", ""),
			},
		},
		WebhookURL:       getEnv("WEBHOOK_URL", ""),
		PollInterval:     getMs("POLL_INTERVAL_MS", 500),
		RPCTimeout:       getMs("RPC_TIMEOUT_MS", 2000),
		StopDetectWindow: getSec("STOP_WINDOW_SEC", 5),
	}
}
