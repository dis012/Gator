package main

import (
	"database/sql"
	"log"
	"os"

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
