package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/juaniten/gator/internal/database"
)

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

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		log.Fatalf("Error getting users from database: %v\n", err)
	}

	for i, user := range users {
		if user == s.config.CurrentUserName {
			users[i] += " (current)"
		}
		fmt.Printf("* %s\n", users[i])
	}

	return nil
}
