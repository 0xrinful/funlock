package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/0xrinful/funlock/internal/models"
)

func (app *application) openCommand() *cli.Command {
	return &cli.Command{
		Name:      "open",
		Usage:     "Open a locked app and track usage duration",
		ArgsUsage: "[app-name]",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "index",
				Usage:   "1-based position of the actual app name in the command (e.g. '3' in 'env VAR=1 appname')",
				Value:   1,
				Aliases: []string{"i"},
			},
		},
		Action: app.openAction,
	}
}

func (app *application) openAction(c *cli.Context) error {
	index := c.Int("index")

	args := c.Args().Slice()

	if len(args) < 1 {
		return cli.Exit(
			fmt.Sprintf(
				"%sUsage: funlock open [app] %s\n%sError: app cmd is required.%s",
				Yellow, Reset, Red, Reset,
			), 1,
		)
	}

	if index < 1 || index > len(args) {
		return cli.Exit(
			fmt.Sprintf(
				"%sInvalid --index: got %d, but only %d arguments provided%s",
				Red, index, len(args), Reset),
			1,
		)
	}

	state, err := app.models.State.Get()
	if err != nil {
		return fmt.Errorf("failed to retrieve user state: %w", err)
	}

	if state.XpBalance <= 0 {
		return cli.Exit(
			fmt.Sprintf(
				"%sYou have %s%d XP. %sEarn more XP before unlocking an app.%s",
				Yellow, Green, state.XpBalance, Yellow, Reset,
			), 1,
		)
	}

	if state.IsSessionRunning() {
		return cli.Exit(
			Yellow+"A session is already running. Finish or cancel it before opening another app."+Reset,
			1,
		)
	}

	cmd := exec.Command(strings.Join(args, " "))
	start := time.Now()

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("process error: %w", err)
	}

	end := time.Now()
	duration := end.Sub(start)
	spentXP := app.calculateSpentXP(duration)
	state.XpBalance -= spentXP
	err = app.models.State.Update(state)
	if err != nil {
		return fmt.Errorf("failed to update user state after starting session: %w", err)
	}

	appName := args[index-1]
	session := models.FunSession{
		App:       appName,
		StartTime: start,
		EndTime:   end,
	}
	err = app.models.FunSessions.Insert(&session)
	if err != nil {
		return fmt.Errorf("failed to save fun session [%s]: %w", appName, err)
	}

	msg := fmt.Sprintf(
		"%s✗ Finished fun session with app [%s].%s\n%s⏱ Duration: %s → %d XP spent (%d XP remaining).%s",
		Red,
		session.App,
		Reset,
		Yellow,
		durationStr(duration),
		spentXP,
		state.XpBalance,
		Reset,
	)
	fmt.Println(msg)
	return nil
}
