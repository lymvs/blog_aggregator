package config

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/lymvs/blog_aggregator/internal/database"
)

type Command struct {
	Name      string
	ArgsSlice []string
}

type Commands struct {
	HandlersMap map[string]func(*State, Command) error
}

func (c *Commands) Run(s *State, cmd Command) error {
	if r, found := c.HandlersMap[cmd.Name]; found {
		if err := r(s, cmd); err != nil {
			return err
		}
		return nil
	}
	return errors.New("command not found")
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.HandlersMap[name] = f
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.ArgsSlice) == 0 {
		return errors.New("no arguments passed")
	}

	if _, err := s.Db.GetUser(context.Background(), cmd.ArgsSlice[0]); err != nil {
		os.Exit(1)
	}

	s.Cfg.CurrentUserName = cmd.ArgsSlice[0]

	err := s.Cfg.SetUser()
	if err != nil {
		return err
	}

	fmt.Println("User has been set successfully")
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.ArgsSlice) == 0 {
		return errors.New("no arguments passed")
	}

	args := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.ArgsSlice[0],
	}
	newUser, err := s.Db.CreateUser(context.Background(), args)
	if err != nil {
		os.Exit(1)
	}

	s.Cfg.CurrentUserName = newUser.Name
	err = s.Cfg.SetUser()
	if err != nil {
		return err
	}
	fmt.Println("User registered sucessfully.")
	return nil
}

func Reset(s *State, cmd Command) error {
	if err := s.Db.DeleteUsers(context.Background()); err != nil {
		fmt.Println("Users couldn't be deleted")
		os.Exit(1)
	}
	fmt.Println("Users deleted sucessfully")
	return nil
}
