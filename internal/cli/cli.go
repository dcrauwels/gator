package cli

import (
	"fmt"

	"github.com/dcrauwels/gator/internal/config"
)

type State struct {
	Config *config.Config
}

type Command struct {
	Name      string
	Arguments []string
}

func HandlerLogin(s *State, cmd Command) error {
	// argument sanity check
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("login takes one argument")
	}

	// set username in state
	err := s.Config.SetUser(cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("error setting user: %w", err)
	}

	fmt.Println("User has been set.")

	return nil
}

type Commands struct {
	Call map[string]func(*State, Command) error
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.Call[name] = f
}

func (c *Commands) Run(s *State, cmd Command) error {
	f, ok := c.Call[cmd.Name]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmd.Name)
	}
	return f(s, cmd)
}
