package main

import (
	"context"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"net/url"
	"os"
	"strconv"
	"time"
)

const (
	Location    = "Asia/Tokyo"
	BotUserName = "git commit checker"
)

func init() {
	loc, err := time.LoadLocation(Location)
	if err != nil {
		loc = time.FixedZone(Location, 9*60*60)
	}
	time.Local = loc
}

func main() {
	// envの読み込み
	err := godotenv.Load("./.env")
	if err != nil {
		fmt.Println("[ERROR] env load error")
		os.Exit(1)
	}

	fmt.Println("process start")
	fmt.Println("run environment is " + os.Getenv("ENV"))
	numOfActivity := getGithubUserInfo()

	fmt.Println("num of activity　" + strconv.Itoa(numOfActivity))
	postMessage := getPostComment(numOfActivity)
	postTwitter(postMessage)
	postSlack(postMessage)
	fmt.Println("[INFO] finish commit check")
}

// 当日のアクティビティを取得する（privateのコミット、アクティビティも含む）
func getGithubUserInfo() int {
	var accessToken string = os.Getenv("GITHUB_PERSONAL_ACCESS_KEY")
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	client := github.NewClient(oauth2.NewClient(ctx, ts))
	// privateなアクティビティも取得するので、keyにはprivateなコミットにもアクセスできるだけの権限を与えておくこと
	res, _, err := client.Activity.ListEventsPerformedByUser(ctx, "lmimsra", false, &github.ListOptions{
		Page:    0,
		PerPage: 0,
	})

	//resCommits, _, _ := client.Repositories.ListCommits(ctx, "lmimsra", "privante", &github.CommitsListOptions{
	//	Author: "lmimsra",
	//	ListOptions: github.ListOptions{
	//		Page:    0,
	//		PerPage: 0,
	//	},
	//})

	// 実行日の0時0分0秒を作成
	now := time.Now()
	unixTime := now.Unix()
	_, offset := now.Zone()
	today := time.Unix((unixTime/86400)*86400-int64(offset), 0)
	year, month, day := today.Date()
	fmt.Printf("today is %v年 %v月 %v日\n", year, int(month), day)

	var todayActivity []*github.Event
	if err == nil {
		for i := range res {
			fmt.Println("repoName: " + res[i].Repo.GetName() + " ActivityTime: " + res[i].CreatedAt.Local().String())

			if res[i].CreatedAt.Local().After(today) {
				todayActivity = append(todayActivity, res[i])
			}
		}
	} else {
		fmt.Println(err)
		fmt.Println("[ERROR] commit get failed from github.....")
		os.Exit(1)
	}
	fmt.Println("today Activity length is " + strconv.Itoa(len(todayActivity)))
	fmt.Println("today Activity")
	for i := range todayActivity {
		fmt.Println("repoName: " + todayActivity[i].Repo.GetName() + "Create: " + res[i].CreatedAt.Local().String())
	}
	return len(todayActivity)
}

// アクティビティ数をTwitterに投稿
func postTwitter(postMessage string) {
	apiKey := os.Getenv("TWITTER_API_KEY")
	apiSecret := os.Getenv("TWITTER_API_SECRET")
	accessToken := os.Getenv("TWITTER_ACCESS_TOKEN")
	accessTokenSecret := os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")
	api := anaconda.NewTwitterApiWithCredentials(accessToken, accessTokenSecret, apiKey, apiSecret)
	targetTwitterAccount := "@" + os.Getenv("TWITTER_TARGET_USER_ID")
	message := BotUserName + " posted message\n" + targetTwitterAccount + "\n" + postMessage

	tweet, err := api.PostTweet(message, url.Values{})
	if err != nil {
		fmt.Println("[ERROR] Slack post fail")
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println("[INFO] tweet success!  post ID:" + tweet.IdStr)
	}
}

// アクティビティ数をSlackに投稿
func postSlack(postMessage string) {
	webHookURL := os.Getenv("SLACK_WEB_HOOK_URL")
	field1 := slack.Field{Title: "コミット数のお知らせ", Value: postMessage}
	attachment := slack.Attachment{}
	attachment.AddField(field1)
	color := "good"
	attachment.Color = &color
	targetSlackAccount := fmt.Sprintf("<@%s>", os.Getenv("SLACK_TARGET_USER_ID"))

	payload := slack.Payload{
		Parse:       "",
		Username:    BotUserName,
		IconUrl:     "",
		IconEmoji:   "",
		Channel:     "",
		Text:        targetSlackAccount,
		LinkNames:   "",
		Attachments: []slack.Attachment{attachment},
		UnfurlLinks: false,
		UnfurlMedia: false,
		Markdown:    false,
	}
	err := slack.Send(webHookURL, "", payload)
	if err != nil {
		fmt.Println("[ERROR] Slack post fail")
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println("[INFO] slack post success!")
	}
}

// 送るメッセージ内容をアクティビティ数で変更
func getPostComment(numOfActivity int) string {
	var message string
	switch {
	case numOfActivity == 0:
		message = "今日のコミット数0！？\nまずいですよ！！！"
	case numOfActivity >= 1 && numOfActivity < 5:
		message = fmt.Sprintf("今日はコミット%d件\nまぁまぁやってはいるようですね", numOfActivity)
	case numOfActivity >= 5:
		message = fmt.Sprintf("今日はコミット%d件も！\nよろしい！", numOfActivity)
	}
	return message
}
