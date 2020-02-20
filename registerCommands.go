package main

type positionData struct {
	title       string
	value       float64
	dayIncrease float64
	percentNet  float64
}

func registerCommands() {
	registerAlternate("h", "help")
	registerCommand("help", cmdHelp)

	registerAlternate("bal", "balance")
	registerAlternate("b", "balance")
	registerAlternate("money", "balance")
	registerAlternate("$", "balance")
	registerCommand("balance", cmdBalance)

	registerAlternate("stocks", "list")
	registerAlternate("shares", "list")
	registerAlternate("ls", "list")
	registerAlternate("positions", "list")
	registerCommand("list", cmdList)

	registerAlternate("px", "price")
	registerCommand("price", cmdPrice)

	registerCommand("reset", cmdReset)

	registerCommand("cancel", cmdCancel)

	registerCommand("stop", cmdStop)

	registerCommand("buy", cmdBuy)

	registerCommand("sell", cmdSell)

	registerAlternate("lb", "top")
	registerAlternate("leaders", "top")
	registerAlternate("leaderboard", "top")
	registerAlternate("leaderboards", "top")
	registerCommand("top", cmdTop)

	registerCommand("newseason", cmdNewSeason)

	registerCommand("primary", cmdPrimary)
}
