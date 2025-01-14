package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/juaniten/gator/internal/config"
	"github.com/juaniten/gator/internal/database"

	_ "github.com/lib/pq"
)

type state struct {
	config *config.Config
	db     *database.Queries
}

func main() {
	// Read JSON config file and load into state
	configuration, err := config.Read()
	if err != nil {
		log.Fatalf("error reading gator configuration: %v", err)
	}

	// Load DB
	dbURL := configuration.DbURL
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error opening postgres database: %v", err)
	}
	dbQueries := database.New(db)

	state_pointer := &state{
		config: &configuration,
		db:     dbQueries,
	}

	// Initialize commands
	comm := commands{
		handlers: map[string]func(*state, command) error{},
	}
	comm.register("login", handlerLogin)
	comm.register("register", handlerRegister)
	comm.register("reset", handlerReset)
	comm.register("users", handlerUsers)
	comm.register("agg", handlerAgg)
	comm.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	comm.register("feeds", handlerFeeds)
	comm.register("follow", middlewareLoggedIn(handlerFollow))
	comm.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	comm.register("following", middlewareLoggedIn(handlerFollowing))
	comm.register("browse", middlewareLoggedIn(handlerBrowse))

	// Process arguments
	args := os.Args
	if len(args) < 2 {
		log.Fatal("command name needed")
	}
	// Execute command
	newCommand := command{
		Name:      args[1],
		Arguments: args[2:],
	}
	err = comm.run(state_pointer, newCommand)
	if err != nil {
		fmt.Printf("error running command: %v\n", err)
		os.Exit(1)
	}
}
