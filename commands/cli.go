package commands

import "log/slog"

type CLI struct {
	Build     Build      `cmd:"" help:"builds database by downloading PBF"`
	Convert   Convert    `cmd:"" help:"convert PBF to sqlite"`
	Generate  Generate   `cmd:"" help:"generate a query to sql (helpful for piping into sqlite)"`
	Integrity Integrity  `cmd:"" help:"check integrity between uncompressed and compressed databases"`
	Runtime   Runtime    `cmd:"" help:"run javascript files locally"`
	Server    Server     `cmd:"" help:"start the API server"`
	LogLevel  slog.Level `help:"Set the log level (debug, info, warn, error)" default:"info"`
}
