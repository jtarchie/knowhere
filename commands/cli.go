package commands

type CLI struct {
	Build     Build     `cmd:""`
	Convert   Convert   `cmd:""`
	Integrity Integrity `cmd:""`
	Generate  Generate  `cmd:""`
	Runtime   Runtime   `cmd:""`
	Server    Server    `cmd:""`
}
