package admob

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/google/subcommands"
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

func (p *NoticeSummaryCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	for _, arg := range f.Args() {
		fmt.Println(os.Getenv("ADMOB_PUBLISHER_ID"))
		fmt.Printf("%s ", arg)
	}
	fmt.Println()
	return subcommands.ExitSuccess
}
