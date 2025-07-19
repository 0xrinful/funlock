package main

import (
	"fmt"
	"time"

	"github.com/urfave/cli/v2"
)

func (app *application) finishSessionCommand() *cli.Command {
	return &cli.Command{
		Name:      "finish",
		Usage:     "Finish the currently running session",
		ArgsUsage: "",
		Action:    app.finishSessionAction,
	}
}

func (app *application) finishSessionAction(c *cli.Context) error {
	state, err := app.models.State.Get()
	if err != nil {
		return fmt.Errorf("failed to retrieve user state: %w", err)
	}

	if !state.IsSessionRunning() {
		msg := fmt.Sprintf(
			"%sNo active session to finish.%s\n%sUse `funlock start` to start a session first.%s",
			Yellow,
			Reset,
			Red,
			Reset,
		)
		return cli.Exit(msg, 1)
	}

	err = app.models.WorkSessions.FinishSession(*state.CurrentSessionID)
	if err != nil {
		return fmt.Errorf("failed to finish the current session: %w", err)
	}

	session, err := app.models.WorkSessions.GetByID(*state.CurrentSessionID)
	if err != nil {
		return fmt.Errorf("failed to fetch the finished session: %w", err)
	}
	duration := session.EndTime.Sub(session.StartTime)
	earnedXP := app.calculateXP(duration)

	state.CurrentSessionID = nil
	state.XpBalance += earnedXP
	err = app.models.State.Update(state)
	if err != nil {
		return fmt.Errorf("failed to update user state after finishing session: %w", err)
	}

	msg := fmt.Sprintf(
		"%s✓ Finished session with tag [%s].%s\n%s⏱ Duration: %s → %d XP earned (%s of fun unlocked).%s",
		Green,
		session.Tag,
		Reset,
		Yellow,
		durationStr(duration),
		earnedXP,
		durationStr(time.Duration(earnedXP)*time.Second),
		Reset,
	)

	fmt.Println(msg)

	return nil
}
