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

// define struct to hold available commands
type Commands struct {
	Call map[string]func(*State, Command) error
}

// this struct implements a function to register a new command to its map
func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.Call[name] = f
}

// this struct implements a function to run a previously registered command from its map
func (c *Commands) Run(s *State, cmd Command) error {
	f, ok := c.Call[cmd.Name]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmd.Name)
	}
	return f(s, cmd)
}
