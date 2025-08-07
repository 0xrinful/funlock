package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/0xrinful/funlock/internal/models"
)

const (
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Red    = "\033[31m"
	Reset  = "\033[0m"
	pipe   = Green + "|" + Reset
)

func durationStr(duration time.Duration) string {
	duration = duration.Round(time.Second)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	result := ""
	if hours > 0 {
		result += fmt.Sprintf("%dh ", hours)
	}
	if minutes > 0 || hours > 0 {
		result += fmt.Sprintf("%dm ", minutes)
	}
	result += fmt.Sprintf("%ds", seconds)

	return result
}

func (app *application) calculateEarnedXP(duration time.Duration) int64 {
	return int64(duration.Seconds() * app.config.XPFactor)
}

func (app *application) calculateSpentXP(duration time.Duration) int64 {
	return int64(duration.Seconds())
}

func (app *application) printWorkSessionsTable(sessions []*models.WorkSession) {
	line := strings.Repeat("-", 93)
	fmt.Println(green(line))
	fmt.Printf(green("|  %-3s |  %-15s |  %-16s |  %-16s |  %-11s |  %-7s |\n"),
		"ID", "Tag", "Start Time", "End Time", "Duration", "XP")
	fmt.Println(green(line))

	for _, s := range sessions {
		duration := s.EndTime.Sub(s.StartTime)
		xp := app.calculateEarnedXP(duration)
		fmt.Printf(
			pipe+"  %-3d "+pipe+"  "+yellow(
				"%-15s",
			)+" "+pipe+"  %-16s "+pipe+"  %-16s "+pipe+"  "+yellow(
				"%-11s",
			)+" "+pipe+"  "+red(
				"%-7d",
			)+" "+pipe+"\n",
			s.ID,
			s.Tag,
			s.StartTime.Format("02-01-2006 15:04"),
			s.EndTime.Format("02-01-2006 15:04"),
			durationStr(duration),
			xp,
		)
		fmt.Println(green(line))
	}
}

func (app *application) printFunSessionsTable(sessions []*models.FunSession) {
	line := strings.Repeat("-", 93)
	fmt.Println(green(line))
	fmt.Printf(green("|  %-3s |  %-15s |  %-16s |  %-16s |  %-11s |  %-7s |\n"),
		"ID", "App", "Start Time", "End Time", "Duration", "XP")
	fmt.Println(green(line))

	for _, s := range sessions {
		duration := s.EndTime.Sub(s.StartTime)
		xp := app.calculateSpentXP(duration)
		fmt.Printf(
			pipe+"  %-3d "+pipe+"  "+yellow(
				"%-15s",
			)+" "+pipe+"  %-16s "+pipe+"  %-16s "+pipe+"  "+yellow(
				"%-11s",
			)+" "+pipe+"  "+red(
				"%-7d",
			)+" "+pipe+"\n",
			s.ID,
			s.App,
			s.StartTime.Format("02-01-2006 15:04"),
			s.EndTime.Format("02-01-2006 15:04"),
			durationStr(duration),
			xp,
		)
		fmt.Println(green(line))
	}
}

func (app *application) printWorkTagSummaries(summaries []*models.WorkTagSummary) {
	line := strings.Repeat("-", 60)
	fmt.Println(green(line))
	fmt.Printf(green("|  %-20s |  %-20s |  %-7s |\n"), "Tag", "Total Duration", "XP")
	fmt.Println(green(line))

	for _, s := range summaries {
		xp := app.calculateEarnedXP(s.Duration)
		fmt.Printf(
			pipe+"  "+yellow(
				"%-20s",
			)+" "+pipe+"  "+yellow(
				"%-20s",
			)+" "+pipe+"  "+red(
				"%-7d",
			)+" "+pipe+"\n",
			s.Tag,
			durationStr(s.Duration),
			xp,
		)
		fmt.Println(green(line))
	}
}

func (app *application) printFunAppSummaries(summaries []*models.FunAppSummary) {
	line := strings.Repeat("-", 60)
	fmt.Println(green(line))
	fmt.Printf(green("|  %-20s |  %-20s |  %-7s |\n"), "App", "Total Duration", "XP")
	fmt.Println(green(line))

	for _, s := range summaries {
		xp := app.calculateEarnedXP(s.Duration)
		fmt.Printf(
			pipe+"  "+yellow(
				"%-20s",
			)+" "+pipe+"  "+yellow(
				"%-20s",
			)+" "+pipe+"  "+red("%-7d")+" "+pipe+"\n",
			s.App,
			durationStr(s.Duration),
			xp,
		)
		fmt.Println(green(line))
	}
}

func (app *application) printWeeklyWorkStats(stats map[time.Weekday]time.Duration) {
	line := strings.Repeat("-", 50)
	fmt.Println(green(line))
	fmt.Printf(green("|  %-16s |  %-25s |\n"), "Day", "Total Work Duration")
	fmt.Println(green(line))

	days := []time.Weekday{
		time.Sunday, time.Monday, time.Tuesday,
		time.Wednesday, time.Thursday, time.Friday, time.Saturday,
	}

	for _, day := range days {
		dur := stats[day]
		durStr := durationStr(dur)
		if dur == 0 {
			durStr = "-"
		}
		fmt.Printf(
			pipe+"  "+yellow("%-16s")+" "+pipe+"  "+yellow("%-25s")+" "+pipe+"\n",
			day.String(),
			durStr,
		)
		fmt.Println(green(line))
	}
}

func containApp(apps []*models.LockedApp, appName string) bool {
	for _, app := range apps {
		if app.Name == appName {
			return true
		}
	}
	return false
}

func green(str string) string {
	return Green + str + Reset
}

func yellow(str string) string {
	return Yellow + str + Reset
}

func red(str string) string {
	return Red + str + Reset
}

func xpToDuratoinStr(xp int64) string {
	sign := ""
	if xp < 0 {
		sign = "-"
		xp = xp * -1
	}
	return sign + durationStr(time.Duration(xp)*time.Second)
}
