package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/dcrauwels/gator/internal/config"
	"github.com/dcrauwels/gator/internal/database"
	"github.com/google/uuid"
)

type State struct {
	Config *config.Config
	Db     *database.Queries
}

type Command struct {
	Name      string
	Arguments []string
}

// set database user
func HandlerLogin(s *State, cmd Command) error {
	// argument sanity check
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("login takes exactly one argument")
	}

	name := cmd.Arguments[0]

	// check if user in db
	ctx := context.Background()
	_, err := s.Db.GetUser(ctx, name)
	if err != nil {
		return fmt.Errorf("user is not registered: %w", err)
	}

	// set username in state
	err = s.Config.SetUser(name)
	if err != nil {
		return fmt.Errorf("error setting user: %w", err)
	}

	fmt.Printf("User has been set to '%s'.\n", name)

	return nil
}

// register database user
func HandlerRegister(s *State, cmd Command) error {
	// argument sanity check
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("register takes exactly one argument")
	}

	// is this even correct? I have no clue, just taken from PSQL docs
	ctx := context.Background()

	// params struct
	name := cmd.Arguments[0]
	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}

	// execute query: call database.CreateUser()
	insertedUser, err := s.Db.CreateUser(ctx, params)
	if err != nil {
		return fmt.Errorf("error inserting user into database: %w", err)
	}

	// set user
	err = s.Config.SetUser(name)
	if err != nil {
		return fmt.Errorf("error setting user: %w", err)
	}

	// log to term
	fmt.Println("User created:")
	fmt.Println(insertedUser)

	return nil
}

func HandlerReset(s *State, cmd Command) error {
	// argument sanity check
	if len(cmd.Arguments) > 0 {
		return fmt.Errorf("register takes exactly zero arguments")
	}

	// execute query: call database.ResetUser()
	ctx := context.Background()
	if err := s.Db.ResetUser(ctx); err != nil {
		return fmt.Errorf("error resetting table: %w", err)
	}

	return nil
}

func HandlerUsers(s *State, cmd Command) error {
	// argument sanity check
	if len(cmd.Arguments) > 0 {
		return fmt.Errorf("register takes exactly zero arguments")
	}

	// execute query: call database.GetUsers()
	ctx := context.Background()
	users, err := s.Db.GetUsers(ctx)
	if err != nil {
		return fmt.Errorf("error getting users: %w", err)
	}

	var uName string

	for _, u := range users {
		uName = u.Name
		if uName == s.Config.CurrentUserName {
			uName += " (current)"
		}
		fmt.Printf("* %s\n", uName)
	}

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
