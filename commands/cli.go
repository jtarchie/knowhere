package commands

type CLI struct {
	Build     Build     `cmd:""`
	Convert   Convert   `cmd:""`
	Integrity Integrity `cmd:""`
	Query     Generate  `cmd:""`
	Server    Server    `cmd:""`
}
