package main

import (
	"log/slog"
	"os"

	"github.com/alecthomas/kong"
	"github.com/jtarchie/knowhere/commands"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, nil)))

	cli := &commands.CLI{}
	ctx := kong.Parse(cli)
	// Call the Run() method of the selected parsed command.
	err := ctx.Run(os.Stdout)
	ctx.FatalIfErrorf(err)
}
