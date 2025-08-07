package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

func (app *application) showCommand() *cli.Command {
	return &cli.Command{
		Name:      "show",
		Usage:     "Show last [n] fun or work sessions",
		ArgsUsage: "[fun|work] [n]",
		Action:    app.showAction,
	}
}

func (app *application) showAction(c *cli.Context) error {
	mode := c.Args().First()
	if mode == "" {
		return cli.Exit(
			fmt.Sprintf(
				"%sUsage: funlock show [fun|work|state|tags|apps|stats] [count]%s\n%sError: show mode is required.%s",
				Yellow,
				Reset,
				Red,
				Reset,
			), 1,
		)
	}

	countStr := "10"
	if c.NArg() >= 2 {
		countStr = c.Args().Get(1)
	}
	count, err := strconv.Atoi(countStr)
	if err != nil || count <= 0 {
		return cli.Exit(fmt.Sprintf("%sInvalid count: %s%s\n", Red, countStr, Reset), 1)
	}

	switch mode {
	case "work":
		sessions, err := app.models.WorkSessions.GetLastN(count)
		if err != nil {
			return fmt.Errorf("failed to fetch work sessions: %w", err)
		}
		app.printWorkSessionsTable(sessions)
	case "fun":
		sessions, err := app.models.FunSessions.GetLastN(count)
		if err != nil {
			return fmt.Errorf("failed to fetch fun sessions: %w", err)
		}
		app.printFunSessionsTable(sessions)
	case "state":
		state, err := app.models.State.Get()
		if err != nil {
			return fmt.Errorf("failed to retrieve user state: %w", err)
		}
		duration, err := app.models.WorkSessions.TotalWorkDuration()
		if err != nil {
			return fmt.Errorf("failed to retrieve work duration: %w", err)
		}
		line := strings.Repeat("-", 35)
		fmt.Printf("%s%s%s\n", Green, line, Reset)
		fmt.Printf(
			"%s★ Current XP: %s%d XP%s (%s)%s\n",
			Yellow,
			Green,
			state.XpBalance,
			Yellow,
			xpToDuratoinStr(state.XpBalance),
			Reset,
		)
		fmt.Printf(
			"%s⏱ Total Work Duration: %s%s %s\n",
			Yellow,
			Green,
			durationStr(duration),
			Reset,
		)
		fmt.Printf("%s%s%s\n", Green, line, Reset)
	case "tags":
		tagSummaries, err := app.models.WorkSessions.GetWorkTimeByTag(count)
		if err != nil {
			return fmt.Errorf("failed to fetch tag summaries: %w", err)
		}
		app.printWorkTagSummaries(tagSummaries)
	case "apps":
		appSummaries, err := app.models.FunSessions.GetFunTimeByApp(count)
		if err != nil {
			return fmt.Errorf("failed to fetch app summaries: %w", err)
		}
		app.printFunAppSummaries(appSummaries)
	case "stats":
		stats, err := app.models.WorkSessions.GetWeeklyWorkStats()
		if err != nil {
			return fmt.Errorf("failed to fetch week stats: %w", err)
		}
		app.printWeeklyWorkStats(stats)
	default:
		return cli.Exit(fmt.Sprintf("%sUnknown Mode type: %s%s\n", Red, mode, Reset), 1)
	}
	return nil
}
