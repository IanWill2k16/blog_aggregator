package main

import (
	"fmt"

	"github.com/IanWill2k16/blog_aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	cfg.SetUser("ian")
	new_cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(new_cfg)
}
