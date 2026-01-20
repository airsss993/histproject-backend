package app

import (
	"fmt"
	"log"

	"github.com/airsss993/histproject-backend/internal/config"
	"github.com/airsss993/histproject-backend/pkg/db"
)

func Run() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal("failed to init config. ", err)
	}
	_ = db.ConnDB(cfg)

	fmt.Println("Hello, World!")
}
