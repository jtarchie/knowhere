package commands

import "log/slog"

type CLI struct {
	Build    Build      `cmd:""`
	Convert  Convert    `cmd:""`
	Generate Generate   `cmd:""`
	Runtime  Runtime    `cmd:""`
	Server   Server     `cmd:""`
	LogLevel slog.Level `help:"Set the log level (debug, info, warn, error)" default:"info"`
}
