package config

import (
	"errors"
	"fmt"
)

type Command struct {
	Name 		string
	ArgsSlice	[]string
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

	s.Cfg.CurrentUserName = cmd.ArgsSlice[0]

	err := s.Cfg.SetUser()
	if err != nil {
		return err
	}

	fmt.Println("User has been set successfully")
	return nil
}