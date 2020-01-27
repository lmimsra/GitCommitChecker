package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"os"
	"time"
)

const location = "Asia/Tokyo"

func init() {
	loc, err := time.LoadLocation(location)
	if err != nil {
		loc = time.FixedZone(location, 9*60*60)
	}
	time.Local = loc
}

func main() {
	// envの読み込み
	err := godotenv.Load("./.env")
	if err != nil {
		fmt.Println("env load error")
		os.Exit(1)
	}

	fmt.Println("process start")
	fmt.Println("run environment is " + os.Getenv("ENV"))
	getGithubUserInfo()
}

func getGithubUserInfo() {
	var accessToken string = os.Getenv("GITHUB_PERSONAL_ACCESS_KEY")
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	client := github.NewClient(oauth2.NewClient(ctx, ts))
	res, _, err := client.Activity.ListEventsPerformedByUser(ctx, "lmimsra", false, &github.ListOptions{
		Page:    0,
		PerPage: 0,
	})
	now := time.Now()
	unixTime := now.Unix()
	_, offset := now.Zone()
	today := time.Unix((unixTime/86400)*86400-int64(offset), 0)
	fmt.Println("toDay is " + today.String())
	var todayActivity []*github.Event
	if err == nil {
		for i := range res {
			if res[i].CreatedAt.Local().After(today) {
				todayActivity = append(todayActivity, res[i])
			}
		}
	} else {
		fmt.Println(err)
	}

	for i := range todayActivity {
		fmt.Println("today Activity")
		fmt.Println("repoName: " + todayActivity[i].Repo.GetName() + "Create: " + res[i].CreatedAt.Local().String())
	}
	//fmt.Println(res[1])
}
