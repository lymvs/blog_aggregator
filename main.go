package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/lymvs/blog_aggregator/internal/config"
	"github.com/lymvs/blog_aggregator/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(errors.New("error reading the file"))
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		fmt.Println(errors.New("error connecting to database"))
		os.Exit(1)
	}
	dbQueries := database.New(db)

	state := config.State{
		Db:  dbQueries,
		Cfg: &cfg,
	}

	commands := config.Commands{
		HandlersMap: make(map[string]func(*config.State, config.Command) error),
	}

	commands.Register("login", config.HandlerLogin)
	commands.Register("register", config.HandlerRegister)
	commands.Register("reset", config.Reset)
	commands.Register("users", config.Users)
	commands.Register("agg", config.Agg)

	args := os.Args[1:] //ignore first argument, that is the program name

	if len(args) < 1 {
		fmt.Println("too few arguments passed")
		os.Exit(1)
	}

	command := config.Command{
		Name:      args[0],
		ArgsSlice: args[1:],
	}

	if err := commands.Run(&state, command); err != nil {
		fmt.Printf("command execution failed: %v\n", err)
		os.Exit(1)
	}
}
