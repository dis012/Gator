package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dis012/agreGator/internal/config"
	"github.com/dis012/agreGator/internal/database"

	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading the config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.Db_url)
	if err != nil {
		log.Fatalf("error connecting to the database: %v", err)
	}

	db_queries := database.New(db)

	program_state := &state{
		db:  db_queries,
		cfg: &cfg,
	}

	cmds := commands{
		available_commands: map[string]func(*state, command) error{},
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handleRegister)
	cmds.register("reset", handleReset)
	cmds.register("users", handleUsers)
	cmds.register("agg", middlewareFetcs(handleAggregate))
	cmds.register("addfeed", middlewareLoggedIn(handleAddFeed))
	cmds.register("feeds", handleFeeds)
	cmds.register("follow", middlewareLoggedIn(handleFollow))
	cmds.register("following", middlewareLoggedIn(handleFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handleUnfollow))

	args := os.Args

	if len(args) < 2 {
		log.Fatal("not enough arguments provided")
	}

	cmd := command{
		name:      args[1],
		arguments: args[2:],
	}

	err = cmds.run(program_state, cmd)
	if err != nil {
		log.Fatal(err)
	}
}

func middlewareLoggedIn(handler func(*state, command, database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.Current_user_name)
		if err != nil {
			return fmt.Errorf("error getting user from database: %v", err)
		}

		return handler(s, cmd, user)
	}
}

func middlewareFetcs(handler func(*state, command, time.Duration) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		time, _ := time.ParseDuration("2s")
		return handler(s, cmd, time)
	}
}
