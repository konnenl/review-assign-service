package main

import (
	"fmt"
	"github.com/konnen/review-assign-service/internal/config"
	"os"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("failed to load config: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(cfg)
}
