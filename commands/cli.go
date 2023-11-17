package commands

type CLI struct {
	Build  Build  `cmd:""`
	Query  Query  `cmd:""`
	Server Server `cmd:""`
}
