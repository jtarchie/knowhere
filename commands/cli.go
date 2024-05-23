package commands

type CLI struct {
	Build   Build   `cmd:""`
	Convert Convert `cmd:""`
	Query   Query   `cmd:""`
	Server  Server  `cmd:""`
}
