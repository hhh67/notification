package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"
	"github.com/hhh67/notification/admob"
)

func main() {
	subcommands.Register(&admob.NoticeSummaryCmd{}, "aa")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
