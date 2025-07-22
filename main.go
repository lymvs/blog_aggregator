package main

import (
	"fmt"
	"errors"

	"github.com/lymvs/blog_aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(errors.New("error reading the file"))
	}

	err = cfg.SetUser()
	if err != nil {
		fmt.Println(errors.New("error setting the username"))
	}

	cfg_updated, err := config.Read()
	if err != nil {
		fmt.Println(errors.New("error reading the file"))
	}

	fmt.Printf("db_url: %s\n", cfg_updated.DbURL)
	fmt.Printf("current_user_name: %s\n", cfg_updated.CurrentUserName)
}