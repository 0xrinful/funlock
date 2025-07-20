package main

import (
	"fmt"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/0xrinful/funlock/internal/models"
)

func (app *application) startCommand() *cli.Command {
	return &cli.Command{
		Name:      "start",
		Usage:     "Start a new focus session assigned to a tag",
		ArgsUsage: "[tag]",
		Action:    app.startAction,
	}
}

func (app *application) startAction(c *cli.Context) error {
	state, err := app.models.State.Get()
	if err != nil {
		return fmt.Errorf("failed to retrieve user state: %w", err)
	}

	if state.IsSessionRunning() {
		session, err := app.models.WorkSessions.GetByID(*state.CurrentSessionID)
		if err != nil {
			return fmt.Errorf("failed to fetch currently running session: %w", err)
		}
		duration := time.Since(session.StartTime)
		msg := fmt.Sprintf(
			"%sA session with tag [%s] is already running (started %s ago).%s\n%sUse `funlock finish` before starting a new one.%s",
			Yellow,
			session.Tag,
			durationStr(duration),
			Reset,
			Red,
			Reset,
		)
		return cli.Exit(msg, 1)
	}

	tag := c.Args().First()
	if tag == "" {
		return cli.Exit(
			fmt.Sprintf(
				"%sUsage: funlock start [tag]%s\n%sError: tag name is required.%s",
				Yellow, Reset, Red, Reset,
			),
			1,
		)
	}

	session := &models.WorkSession{
		Tag:       tag,
		StartTime: time.Now(),
	}

	id, err := app.models.WorkSessions.Insert(session)
	if err != nil {
		return fmt.Errorf("failed to start session [%s]: %w", tag, err)
	}

	state.CurrentSessionID = &id
	err = app.models.State.Update(state)
	if err != nil {
		return fmt.Errorf("failed to update user state after starting session: %w", err)
	}

	fmt.Printf("%s✓ Started session with tag [%s]%s\n", Green, tag, Reset)
	return nil
}
