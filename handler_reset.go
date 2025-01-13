package main

import (
	"context"
	"fmt"
	"log"
	"os"
)

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		log.Fatalf("Error resetting database: %v\n", err)
	}
	fmt.Println("Database successfuly reseted.")
	os.Exit(0)
	return nil
}
