package main

import (
	"io"
	"log/slog"
	"os"

	"github.com/alecthomas/kong"
	"github.com/jtarchie/knowhere/commands"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, nil)))

	cli := &commands.CLI{}
	var writer io.Writer = os.Stdout

	ctx := kong.Parse(cli, kong.BindTo(writer, (*io.Writer)(nil)))
	// Call the Run() method of the selected parsed command.
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
