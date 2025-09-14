package config

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/lymvs/blog_aggregator/internal/database"
	"github.com/lymvs/blog_aggregator/internal/rss"
)

type Command struct {
	Name      string
	ArgsSlice []string
}

type Commands struct {
	HandlersMap map[string]func(*State, Command) error
}

func (c *Commands) Run(s *State, cmd Command) error {
	if r, found := c.HandlersMap[cmd.Name]; found {
		if err := r(s, cmd); err != nil {
			return err
		}
		return nil
	}
	return errors.New("command not found")
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.HandlersMap[name] = f
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.ArgsSlice) == 0 {
		return errors.New("no arguments passed")
	}

	if _, err := s.Db.GetUser(context.Background(), cmd.ArgsSlice[0]); err != nil {
		os.Exit(1)
	}

	s.Cfg.CurrentUserName = cmd.ArgsSlice[0]

	err := s.Cfg.SetUser()
	if err != nil {
		return err
	}

	fmt.Println("User has been set successfully")
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.ArgsSlice) == 0 {
		return errors.New("no arguments passed")
	}

	args := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.ArgsSlice[0],
	}
	newUser, err := s.Db.CreateUser(context.Background(), args)
	if err != nil {
		os.Exit(1)
	}

	s.Cfg.CurrentUserName = newUser.Name
	err = s.Cfg.SetUser()
	if err != nil {
		return err
	}
	fmt.Println("User registered sucessfully.")
	return nil
}

func Reset(s *State, cmd Command) error {
	if err := s.Db.DeleteUsers(context.Background()); err != nil {
		fmt.Println("Users couldn't be deleted")
		os.Exit(1)
	}
	fmt.Println("Users deleted sucessfully")
	return nil
}

func Users(s *State, cmd Command) error {
	users_list, err := s.Db.GetUsers(context.Background())

	if err != nil {
		return err
	}

	for _, user := range users_list {
		if user.Name == s.Cfg.CurrentUserName {
			fmt.Printf("%s (current)\n", user.Name)
			continue
		}
		fmt.Printf("%s\n", user.Name)
	}
	return nil
}

func Agg(s *State, cmd Command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")

	if err != nil {
		return err
	}

	fmt.Println(feed.Channel.Title)
	fmt.Println(feed.Channel.Link)
	fmt.Println(feed.Channel.Description)
	for i := range feed.Channel.Item {
		fmt.Println(feed.Channel.Item[i].Title)
		fmt.Println(feed.Channel.Item[i].Link)
		fmt.Println(feed.Channel.Item[i].Description)
		fmt.Println(feed.Channel.Item[i].PubDate)
	}

	return nil
}

func fetchFeed(ctx context.Context, feedURL string) (*rss.RSSFeed, error) {
	new_req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		feedURL,
		nil)

	if err != nil {
		return nil, err
	}

	new_req.Header.Set("User-Agent", "gator")

	client := &http.Client{}
	res, err := client.Do(new_req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	rssFeed := rss.RSSFeed{}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if err := xml.Unmarshal(data, &rssFeed); err != nil {
		return nil, err
	}

	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)

	for i := range rssFeed.Channel.Item {
		rssFeed.Channel.Item[i].Title = html.UnescapeString(rssFeed.Channel.Item[i].Title)
		rssFeed.Channel.Item[i].Description = html.UnescapeString(rssFeed.Channel.Item[i].Description)
	}

	return &rssFeed, nil
}

func AddFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.ArgsSlice) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}

	name := cmd.ArgsSlice[0]
	url := cmd.ArgsSlice[1]

	feed, err := s.Db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		Name:      name,
		Url:       url,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed: %w", err)
	}

	if _, err := s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}); err != nil {
		return fmt.Errorf("couldn't create feed row: %w", err)
	}

	fmt.Println("Feed created successfully:")
	printFeed(feed)
	fmt.Println()
	fmt.Println("==================================")

	return nil
}

func printFeed(feed database.Feed) {
	fmt.Printf("* ID: 			%s\n", feed.ID)
	fmt.Printf("* Created:		%v\n", feed.CreatedAt)
	fmt.Printf("* Updated:		%v\n", feed.UpdatedAt)
	fmt.Printf("* Name: 		%s\n", feed.Name)
	fmt.Printf("* URL:			%s\n", feed.Url)
	fmt.Printf("* UserID:		%s\n", feed.UserID)
}

func Feeds(s *State, cmd Command) error {
	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Printf("* Name:			%s\n", feed.Name)
		fmt.Printf("* URL: 			%s\n", feed.Url)
		if feed.Name_2.Valid {
			fmt.Printf("* UserName:		%s\n", feed.Name_2.String)
		} else {
			fmt.Printf("* UserName:		None")
		}
	}

	return nil
}

func Follow(s *State, cmd Command, user database.User) error {
	if len(cmd.ArgsSlice) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	feedID, err := s.Db.GetFeedByURL(context.Background(), cmd.ArgsSlice[0])
	if err != nil {
		return err
	}

	feedRow, err := s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feedID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed row: %w", err)
	}

	fmt.Printf("User %s is now following feed %s\n", feedRow.UserName, feedRow.FeedName)

	return nil
}

func Following(s *State, cmd Command, user database.User) error {
	feedsFollowing, err := s.Db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	for _, feed := range feedsFollowing {
		fmt.Println(feed.FeedName)
	}

	return nil
}

func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(s *State, cmd Command) error {
	return func(s *State, cmd Command) error {
		user, err := s.Db.GetUser(context.Background(), s.Cfg.CurrentUserName)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}

func Unfollow(s *State, cmd Command, user database.User) error {
	feedID, err := s.Db.GetFeedByURL(context.Background(), cmd.ArgsSlice[0])
	if err != nil {
		return err
	}

	if err := s.Db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feedID,
	}); err != nil {
		return err
	}

	return nil
}
