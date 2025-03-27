package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/dcrauwels/gator/internal/cli"
	"github.com/dcrauwels/gator/internal/config"
	"github.com/dcrauwels/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	// set state
	var s cli.State
	cfg, err := config.Read()
	if err != nil {
		panic(err)
	}
	s.Config = &cfg

	db, err := sql.Open("postgres", s.Config.DbUrl)
	if err != nil {
		panic(err)
	}
	dbQueries := database.New(db)
	s.Db = dbQueries

	// init commands
	c := cli.Commands{
		Call: make(map[string]func(*cli.State, cli.Command) error),
	}
	// register commands
	c.Register("login", cli.HandlerLogin)
	c.Register("register", cli.HandlerRegister)
	c.Register("reset", cli.HandlerReset)
	c.Register("users", cli.HandlerUsers)
	c.Register("agg", cli.HandlerAgg)
	c.Register("addfeed", cli.HandlerAddFeed)
	c.Register("feeds", cli.HandlerFeeds)
	c.Register("follow", cli.HandlerFollow)
	c.Register("following", cli.HandlerFollowing)

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

}
