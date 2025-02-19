package main

import (
	"fmt"
	"os"

	"github.com/dcrauwels/gator/internal/cli"
	"github.com/dcrauwels/gator/internal/config"
)

func main() {
	// set state
	var s cli.State
	cfg, err := config.Read()
	if err != nil {
		panic(err)
	}
	s.Config = &cfg

	// init commands
	c := cli.Commands{
		Call: make(map[string]func(*cli.State, cli.Command) error),
	}
	// register commands
	c.Register("login", cli.HandlerLogin)

	// get cli args and sanity check
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Error: no argument called")
		os.Exit(1)
	}

	// parse cli args
	var arguments []string
	if len(args) >= 3 {
		arguments = args[2:]
	}

	// make Command struct and run from Commands
	cmd := cli.Command{
		Name:      args[1],
		Arguments: arguments,
	}

	err = c.Run(&s, cmd)
	if err != nil {
		fmt.Println(fmt.Errorf("error running command: %w", err))
		os.Exit(1)
	} else {
		os.Exit(0)
	}

	/*


		currentUser := "dcrauwels"
		c.SetUser(currentUser)

		c, _ = config.Read()
		fmt.Println(c)
	*/
}
