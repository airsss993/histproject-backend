package queue

import (
	"github.com/airsss993/histproject-backend/internal/config"
	"github.com/hibiken/asynq"
)

type Client struct {
	*asynq.Client
}

func NewClient(cfg config.Redis) *Client {
	return &Client{
		asynq.NewClient(asynq.RedisClientOpt{
			Addr:     cfg.RedisAddress,
			DB:       cfg.RedisDb,
			Password: cfg.RedisPassword,
		}),
	}
}
