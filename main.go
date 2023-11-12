package main

import (
	"log/slog"
	"os"

	"github.com/alecthomas/kong"
	"github.com/jtarchie/knowhere/commands"
)

type CLI struct {
	Build commands.Build `cmd:""`
}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, nil)))

	cli := &CLI{}
	ctx := kong.Parse(cli)
	// Call the Run() method of the selected parsed command.
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
