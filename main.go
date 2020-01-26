package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"os"
)

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
	if nil != err {
		for i := range res {
			fmt.Println("repoName: " + res[i].Repo.GetName() + "Create: " + res[i].CreatedAt.Local().String())
		}
	} else {
		fmt.Println(err)
	}

	//fmt.Println(res[1])
}
