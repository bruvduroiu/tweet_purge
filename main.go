package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Xcod3bughunt3r/Go-Twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/joho/godotenv"
)

func BoolPointer(b bool) *bool {
	return &b
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	config := oauth1.NewConfig(
		os.Getenv("TWITTER_CONSUMER_KEY"),
		os.Getenv("TWITTER_CONSUMER_SECRET"),
	)

	token := oauth1.NewToken(
		os.Getenv("TWITTER_ACCESS_KEY"),
		os.Getenv("TWITTER_ACCESS_SECRET"),
	)

	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	userTimelineParams := &twitter.UserTimelineParams{
		ScreenName:      os.Getenv("TWITTER_USER_HANDLE"),
		IncludeRetweets: BoolPointer(true),
		ExcludeReplies:  BoolPointer(false),
		Count:           200,
	}
	tweets, _, err := client.Timelines.UserTimeline(userTimelineParams)

	if err != nil {
		panic(err)
	}

	f, err := os.Create("tweets.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fErr, err := os.Create("failed_delete.txt")
	if err != nil {
		panic(err)
	}
	defer fErr.Close()

	for ; len(tweets) > 0; tweets, _, _ = client.Timelines.UserTimeline(userTimelineParams) {
		for _, tweet := range tweets {
			fmt.Println("ID=", tweet.IDStr, " Text=", tweet.Text, "Retweet=", tweet.Retweeted)
			if tweet.Retweeted {
				_, _, err = client.Statuses.Unretweet(tweet.ID, &twitter.StatusUnretweetParams{
					ID:       tweet.ID,
					TrimUser: BoolPointer(true),
				})
				if err != nil {
					fErr.WriteString(tweet.IDStr + "," + err.Error())
					log.Fatal("Failed Destroying Tweet ID=", tweet.IDStr)
				}
			} else {
				f.WriteString(tweet.IDStr + "," + tweet.Text + "\n")
				_, _, err = client.Statuses.Destroy(tweet.ID, &twitter.StatusDestroyParams{
					ID:       tweet.ID,
					TrimUser: BoolPointer(true),
				})
				if err != nil {
					fErr.WriteString(tweet.IDStr + "," + err.Error())
					log.Fatal("Failed Unretweeting ID=", tweet.IDStr)
				}
			}
			time.Sleep(1 * time.Second)
		}
	}
}
