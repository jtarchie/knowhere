package commands

import "log/slog"

type CLI struct {
	Build     Build      `cmd:""`
	Convert   Convert    `cmd:""`
	Generate  Generate   `cmd:""`
	Integrity Integrity  `cmd:"" help:"Check integrity between uncompressed and compressed databases"`
	Runtime   Runtime    `cmd:""`
	Server    Server     `cmd:""`
	LogLevel  slog.Level `help:"Set the log level (debug, info, warn, error)" default:"info"`
}
