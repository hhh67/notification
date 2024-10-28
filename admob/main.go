package admob

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/subcommands"
	"github.com/slack-go/slack"
)

type NoticeSummaryCmd struct {
	d bool // 日次
	w bool // 週次
	m bool // 月次
	y bool // 年次
}

func (*NoticeSummaryCmd) Name() string     { return "nss" }
func (*NoticeSummaryCmd) Synopsis() string { return "今日のAdmob売上サマリを通知します" }
func (*NoticeSummaryCmd) Usage() string {
	return `nss [options]:
今日のAdmob売上サマリを通知します
`
}

func (p *NoticeSummaryCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&p.d, "d", false, "日次のサマリを通知します")
	f.BoolVar(&p.w, "w", false, "週次のサマリを通知します")
	f.BoolVar(&p.m, "m", false, "月次のサマリを通知します")
	f.BoolVar(&p.y, "y", false, "年次のサマリを通知します")
}

func (cmd *NoticeSummaryCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	apiUrl := fmt.Sprintf("https://admob.googleapis.com/v1/accounts/pub-%s/networkReport:generate", os.Getenv("ADMOB_PUBLISHER_ID"))

	requestBody, startDate, endDate := MakeRequestBody(cmd)

	result, err := RequestAdmobApi(apiUrl, requestBody)
	if err != nil {
		log.Printf("Admob APIのリクエストに失敗しちゃった！: %v", err)
		return subcommands.ExitFailure
	}

	slackClient := slack.New(os.Getenv("SLACK_API_TOKEN"))
	channel := os.Getenv("SLACK_ADMOB_CHANNEL_ID")
	SendSlackMessage(slackClient, channel, cmd, result, startDate, endDate)

	return subcommands.ExitSuccess
}
