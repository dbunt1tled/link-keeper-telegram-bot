package telegram

const (
	MsgHelp = `I can save and keep your pages, also I can send you random messages.
/start - start bot
/help - this help
/rnd - random message
`
	MsgStart = `
Hello! I'm a bot. I can help you with some tasks.


` + MsgHelp
	MsgUnknown          = `Command not found 🧐`
	MsgNoSavedPages     = `You have no saved pages 🤪`
	MsgSavedPage        = `Saved! 🎉`
	msgAlreadySavedPage = `You already have this page 😉`
)
