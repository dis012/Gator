package main

import (
	"context"
	"fmt"
	"time"

	"github.com/dis012/agreGator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) != 1 {
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

func handleAggregate(s *state, cmd command, time_between_reqs time.Duration) error {
	fmt.Printf("Collecting feeds every %v", time_between_reqs)
	fmt.Println()
	ticker := time.NewTicker(time_between_reqs)
	for ; ; <-ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			return err
		}
	}
}

func handleAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) != 2 {
		return fmt.Errorf("the add feed handler expects a two arguments, the name of the feed and url")
	}

	new_feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
		Url:       cmd.arguments[1],
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating a feed: %v", err)
	}

	fmt.Println("Succesfully created new feed!")
	printFeed(new_feed, user)

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    new_feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed follow: %w", err)
	}

	fmt.Printf("Followed %v", feedFollow.FeedName)

	return nil
}

func handleFeeds(s *state, cmd command) error {
	all_feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error getting feeds: %v", err)
	}

	for _, feed := range all_feeds {
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("error getting user by id: %v", err)
		}
		printFeed(feed, user)
	}

	return nil
}

func printFeed(f database.Feed, u database.User) {
	fmt.Println("=========================================")
	fmt.Printf("Name of feed: %v", f.Name)
	fmt.Println()
	fmt.Printf("Feed url: %v", f.Url)
	fmt.Println()
	fmt.Printf("Creator of feed: %v", u.Name)
	fmt.Println()
	fmt.Println("=========================================")
}

func handleFollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("follow command needs only one parameter and that is url")
	}

	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("error getting the feed: %v", err)
	}

	feed_follow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating feed follow: %v", err)
	}

	fmt.Println("Successfully created feed follow!")
	printFeedFollow(feed_follow)

	return nil
}

func printFeedFollow(f database.CreateFeedFollowRow) {
	fmt.Println("=========================================")
	fmt.Printf("Feed name: %v", f.FeedName)
	fmt.Println()
	fmt.Printf("User name: %v", f.UserName)
	fmt.Println()
	fmt.Println("=========================================")
}

func handleFollowing(s *state, cmd command, user database.User) error {
	feed_follow, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error getting the feed follows: %v", err)
	}

	printAllFeedFollows(feed_follow)

	return nil
}

func printAllFeedFollows(feed_follows []database.GetFeedFollowsForUserRow) {
	fmt.Printf("All feeds that user follows:\n")
	for _, feed := range feed_follows {
		fmt.Println("=========================================")
		fmt.Printf("Feed name: %v", feed.FeedName)
		fmt.Println()
		fmt.Println("=========================================")
	}
}

func handleUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("the unfollow handler expects a single argument, the feed url")
	}

	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("error getting the feed: %v", err)
	}

	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error deleting the feed: %v", err)
	}

	return nil
}
