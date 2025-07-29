package main

import (
	"fmt"
	"errors"
	"os"

	"github.com/lymvs/blog_aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(errors.New("error reading the file"))
	}

	state := config.State{
		Cfg: &cfg,
	}

	commands := config.Commands{
		HandlersMap: make(map[string]func(*config.State, config.Command) error),
	}

	commands.Register("login", config.HandlerLogin)

	args := os.Args[1:] //ignore first argument, that is the program name

	if len(args) < 1 {
		fmt.Println("too few arguments passed")
		os.Exit(1)
	}

	command := config.Command{
		Name: args[0],
		ArgsSlice: args[1:],
	}

	if err := commands.Run(&state, command); err != nil {
		fmt.Printf("command execution failed: %v\n", err)
	}

	cfg_updated, err := config.Read()
	if err != nil {
		fmt.Println(errors.New("error reading the file"))
	}

	fmt.Printf("db_url: %s\n", cfg_updated.DbURL)
	fmt.Printf("current_user_name: %s\n", cfg_updated.CurrentUserName)
}