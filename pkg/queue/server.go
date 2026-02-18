package queue

import (
	"github.com/airsss993/histproject-backend/internal/config"
	"github.com/hibiken/asynq"
)

type Server struct {
	*asynq.Server
}

func NewServer(cfg config.Redis) *Server {
	return &Server{
		asynq.NewServer(asynq.RedisClientOpt{
			Addr:     cfg.RedisAddress,
			DB:       cfg.RedisDb,
			Password: cfg.RedisPassword,
		}, asynq.Config{
			Concurrency: 5,
		}),
	}
}
