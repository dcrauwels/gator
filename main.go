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
	c.Register("login", cli.MwNumArguments(cli.HandlerLogin, 1))
	c.Register("register", cli.MwNumArguments(cli.HandlerRegister, 1))
	c.Register("reset", cli.MwNumArguments(cli.HandlerReset, 0))
	c.Register("users", cli.MwNumArguments(cli.HandlerUsers, 0))
	c.Register("agg", cli.MwNumArguments(cli.HandlerAgg, 1))
	c.Register("addfeed", cli.MwNumArguments(cli.HandlerAddFeed, 2))
	c.Register("feeds", cli.MwNumArguments(cli.HandlerFeeds, 0))
	c.Register("follow", cli.MwNumArguments(cli.HandlerFollow, 1))
	c.Register("following", cli.MwLoggedIn(cli.MwNumArguments(cli.HandlerFollowing, 0)))
	c.Register("unfollow", cli.MwLoggedIn(cli.MwNumArguments(cli.HandlerUnfollow, 1)))

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
