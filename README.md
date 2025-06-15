# agreGator 📰

**agreGator** is a lightweight command-line RSS aggregator written in Go. It connects to a PostgreSQL database, lets you register and manage users, follow feeds, and aggregate content from them in real-time.

---

## 🚀 Prerequisites

To run **agreGator**, you’ll need the following installed:

- **Go** (v1.20 or newer): [Download Go](https://go.dev/dl/)
- **PostgreSQL** (13+): [Install PostgreSQL](https://www.postgresql.org/download/)

---

## 🔧 Installing the `gator` CLI

To install the CLI tool locally, run:

```bash
go install github.com/dis012/agreGator/cmd/gator@latest
```

This will install the binary as gator in your $GOPATH/bin directory. Make sure that directory is in your PATH:
```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

Now you can run gator from anywhere!

## Configuration

Create a configuration file named .gator_config.json in your home directory:
```bash
~/.gator_config.json
```

Example config:

```bash
{
  "db_url": "postgres://username:password@localhost:5432/yourdbname?sslmode=disable",
  "current_user_name": "alice"
}
```

db_url: Your PostgreSQL connection string.
current_user_name: The default user you want to use for logged-in commands.

## Running the program

After installing and setting up your config, just run:
```bash
gator <command> [args...]

Example:
gator register alice
gator login alice
gator addfeed "TechCrunch" "https://techcrunch.com/feed"
gator follow <feed_id>
gator agg

```

## Available Commands
| Command     | Description                                         |
| ----------- | --------------------------------------------------- |
| `register`  | Register a new user                                 |
| `login`     | Log in as a user (sets `current_user_name`)         |
| `addfeed`   | Add a new RSS feed (requires login)                 |
| `feeds`     | List all available feeds                            |
| `follow`    | Follow a feed (requires login)                      |
| `following` | Show feeds the current user is following            |
| `unfollow`  | Unfollow a feed (requires login)                    |
| `agg`       | Aggregate and scrape posts from all feeds           |
| `users`     | List all users                                      |
| `reset`     | Reset user data or config (implementation-specific) |
