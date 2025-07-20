package main

import "github.com/urfave/cli/v2"

func (app *application) commands() []*cli.Command {
	return []*cli.Command{
		app.startCommand(),
		app.finishCommand(),
		app.showCommand(),
		app.openCommand(),
	}
}
