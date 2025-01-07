package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/juaniten/gator/internal/config"
	"github.com/juaniten/gator/internal/database"

	_ "github.com/lib/pq"
)

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

	// Process arguments
	args := os.Args
	if len(args) < 2 {
		log.Fatal("command name needed")
	}
	// Execute command
	newCommand := command{
		name:      args[1],
		arguments: args[2:],
	}
	err = comm.run(state_pointer, newCommand)
	if err != nil {
		fmt.Printf("error running command: %v\n", err)
	}
}

type state struct {
	config *config.Config
	db     *database.Queries
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	if commandFunction, ok := c.handlers[cmd.name]; ok {
		commandFunction(s, cmd)
		return nil
	}
	fmt.Println("command does not exist")
	return nil
}

func handlerLogin(s *state, cmd command) error {

	if len(cmd.arguments) == 0 {
		return errors.New("command called with no arguments")
	}
	username := cmd.arguments[0]

	// Check if username already exists in DB
	user, _ := s.db.GetUser(context.Background(), username)
	if user.Name != username {
		log.Fatalf("Username %s does not exist in DB", username)
	}

	err := s.config.SetUser(username)
	if err != nil {
		return err
	}
	fmt.Printf("Set username '%s'\n", username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("command called with no arguments")
	}
	username := cmd.arguments[0]

	// Check if username already exists in DB
	user, _ := s.db.GetUser(context.Background(), username)
	if user.Name == username {
		log.Fatalf("Username %s already exists in DB", username)
	}

	// Set username
	s.config.SetUser(username)
	args := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username}
	dbUser, err := s.db.CreateUser(context.Background(), args)
	if err != nil {
		log.Fatal("error creating user in DB")
	}

	fmt.Printf("Added user %s\n", username)
	fmt.Printf("User data: %v\n", dbUser)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.Reset(context.Background())
	if err != nil {
		log.Fatalf("Error resetting database: %v\n", err)
	}
	fmt.Println("Database successfuly reseted.")
	os.Exit(0)
	return nil
}
