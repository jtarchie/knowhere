package commands

type CLI struct {
	Convert Convert `cmd:""`
	Query   Query   `cmd:""`
	Server  Server  `cmd:""`
}
