package main

import (
	"context"
	"fmt"
	"time"

	"github.com/dis012/agreGator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("the login handler expects a single argument, the username")
	}

	user, err := s.db.GetUser(context.Background(), cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("error getting user from database: %v", err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("error wiriting to config: %v", err)
	}

	fmt.Println("User has been set!")

	return nil
}

func handleRegister(s *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("the login handler expects a single argument, the username")
	}

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
	})
	if err != nil {
		return fmt.Errorf("error adding user to database: %v", err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("error wiriting to config: %v", err)
	}

	fmt.Println("User has been created!")
	fmt.Printf("User name: %v \n Created at: %v", user.Name, user.CreatedAt)

	return nil
}

func handleReset(s *state, cmd command) error {
	err := s.db.DeletUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error while deleting users: %v", err)
	}

	fmt.Println("Deleted all users from db!")
	return nil
}

func handleUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error while getting all users: %v", err)
	}

	current_user := s.cfg.Current_user_name
	for _, user := range users {
		if user.Name == current_user {
			fmt.Printf("* %v (current)\n", user.Name)
		} else {
			fmt.Printf("* %v\n", user.Name)
		}
	}
	return nil
}
