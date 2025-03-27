package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/dcrauwels/gator/internal/config"
	"github.com/dcrauwels/gator/internal/database"
	"github.com/dcrauwels/gator/internal/rssfeed"
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
	_, err := s.Db.GetUserByName(ctx, name)
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
	fmt.Println("User created succesfully:")
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

func HandlerAgg(s *State, cmd Command) error { // will change later
	// argument sanity check
	if len(cmd.Arguments) > 0 {
		return fmt.Errorf("agg takes exactly zero arguments")
	}

	// get feed url
	feedurl := "https://www.wagslane.dev/index.xml"

	// init context
	ctx := context.Background()

	// fetch feed
	feed, err := rssfeed.FetchFeed(ctx, feedurl)
	if err != nil {
		return fmt.Errorf("error fetching rss: %w", err)
	}

	fmt.Println(feed)

	return nil
}

func HandlerAddFeed(s *State, cmd Command) error {
	// argument sanity check
	if len(cmd.Arguments) != 2 {
		return fmt.Errorf("addfeed takes exactly two arguments")
	}

	// init context
	ctx := context.Background()

	// get current user
	currentUser, err := s.Db.GetUserByName(ctx, s.Config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error getting current user: %w", err)
	}
	currentUserID := currentUser.ID

	// get name, url
	name, url := cmd.Arguments[0], cmd.Arguments[1]

	// run query to add feed to DB CreateFeed
	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    currentUserID,
	}
	createdFeed, err := s.Db.CreateFeed(ctx, params)
	if err != nil {
		return fmt.Errorf("error creating feed: %w", err)
	}

	// log created feed to terminal
	fmt.Println("Feed created succesfully:")
	fmt.Println(createdFeed)

	return nil
}

func HandlerFeeds(s *State, cmd Command) error {
	// argument sanity check
	if len(cmd.Arguments) != 0 {
		return fmt.Errorf("feeds takes exactly zero arguments")
	}

	// init context
	ctx := context.Background()

	// run query to get all rows from db table via GetFeeds
	feeds, err := s.Db.GetFeeds(ctx)
	if err != nil {
		return fmt.Errorf("error getting feeds: %w", err)
	}

	// log feeds to terminal
	fmt.Printf("Number of feeds in table: %d\n", len(feeds))
	for i, f := range feeds {
		// get user name from user id
		u, err := s.Db.GetUserByID(ctx, f.UserID)
		if err != nil {
			return fmt.Errorf("error getting user that created feed: %w", err)
		}
		// log to terminal
		fmt.Printf("\n--Feed %d--\n", i+1)
		fmt.Printf("Name: %s\n", f.Name)
		fmt.Printf("URL: %s\n", f.Url)
		fmt.Printf("Created by: %s\n", u.Name)

	}

	return nil
}

func HandlerFollow(s *State, cmd Command) error {
	// argument sanity check
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("follow takes exactly one argument")
	}

	// get url
	url := cmd.Arguments[0]

	// init context
	ctx := context.Background()

	// check if url in feeds
	feed, err := s.Db.GetFeedByUrl(ctx, url)
	if err != nil {
		return fmt.Errorf("feed is not registered: %w", err)
	}

	// get user ID
	user, err := s.Db.GetUserByName(ctx, s.Config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error checking current user: %w", err)
	}

	// create a new feed follow record for current user
	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	createdFeedFollow, err := s.Db.CreateFeedFollow(ctx, params)
	if err != nil {
		return fmt.Errorf("error creating feed follow: %w", err)
	}

	// print the name of the feed and current user
	fmt.Printf("Followed feed '%s' as user '%s'.\n", createdFeedFollow.FeedName, createdFeedFollow.UserName)
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
