package admob

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/slack-go/slack"
)

func SendSlackMessage(client *slack.Client, channel string, cmd *NoticeSummaryCmd, result []ReportItem, startDate, endDate time.Time) {
	title := fmt.Sprintf("%sのAdmob売上", startDate.Format("2006年1月2日"))
	switch {
	case cmd.w:
		title = "今週のAdmob売上"
	case cmd.m:
		title = "今月のAdmob売上"
	case cmd.y:
		title = fmt.Sprintf("%d年のAdmob売上", startDate.Year())
	}

	var blocks []slack.Block
	blocks = append(blocks, slack.NewSectionBlock(
		nil,
		[]*slack.TextBlockObject{
			slack.NewTextBlockObject(slack.MarkdownType, fmt.Sprintf("*%s*", title+"だよ🔥"), false, false),
		},
		nil,
	))
	var totalEarning float64
	var appEarningFields []*slack.TextBlockObject
	for _, v := range result {
		if v.Footer != nil && v.Footer.MatchingRowCount != "" {
			count, err := strconv.Atoi(v.Footer.MatchingRowCount)
			if err != nil {
				log.Println(err)
				log.Fatalf("予期せぬエラー")
			}
			if count == 0 {
				sendNotFoundMessage(client, channel)
			}
		}
		if v.Row != nil {
			e, err := strconv.Atoi(v.Row.MetricValues.EstimatedEarnings.MicrosValue)
			if err != nil {
				log.Println(err)
				log.Fatalf("予期せぬエラー")
			}
			earning := math.Round(float64(e) / 1000000)
			totalEarning += earning
			appEarningFields = append(appEarningFields, slack.NewTextBlockObject(slack.MarkdownType, fmt.Sprintf("%s: %.1f円\n", v.Row.DimensionValues.App.DisplayLabel, earning), false, false))
		}
	}

	if totalEarning == 0 {
		sendNoEarningMessage(client, channel, title)
		return
	}

	blocks = append(blocks, slack.NewSectionBlock(
		nil,
		appEarningFields,
		nil,
	))

	blocks = append(blocks, slack.NewSectionBlock(
		slack.NewTextBlockObject(slack.MarkdownType, fmt.Sprintf("*合計: %.1f円*💰", totalEarning), false, false),
		nil,
		nil,
	))

	_, _, err := client.PostMessage(channel, slack.MsgOptionBlocks(blocks...))
	if err != nil {
		log.Fatalf("Slackへのメッセージ送信に失敗しちゃった！: %v", err)
	}
}

func sendNotFoundMessage(client *slack.Client, channel string) {
	client.PostMessage(channel, slack.MsgOptionText("Admobのデータが見つからなかった😢", false))
}

func sendNoEarningMessage(client *slack.Client, channel, title string) {
	client.PostMessage(channel, slack.MsgOptionBlocks(
		slack.NewSectionBlock(
			slack.NewTextBlockObject(slack.MarkdownType, "まだ売上がねぇよ...😢", false, false),
			nil,
			nil,
		), slack.NewSectionBlock(
			slack.NewTextBlockObject(slack.MarkdownType, fmt.Sprintf("(%s)", title), false, false),
			nil,
			nil,
		),
	))
}
