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
	fmt.Println(Green + line + Reset)
	fmt.Printf(Green+"|  %-3s |  %-15s |  %-16s |  %-16s |  %-11s |  %-7s |\n"+Reset,
		"ID", "Tag", "Start Time", "End Time", "Duration", "XP")
	fmt.Println(Green + line + Reset)

	for _, s := range sessions {
		duration := s.EndTime.Sub(s.StartTime)
		xp := app.calculateEarnedXP(duration)
		fmt.Printf(
			pipe+"  %-3d "+pipe+"  "+Yellow+"%-15s"+Reset+" "+pipe+"  %-16s "+pipe+"  %-16s "+pipe+"  "+Yellow+"%-11s"+Reset+" "+pipe+"  "+Red+"%-7d"+Reset+" "+pipe+"\n",
			s.ID,
			s.Tag,
			s.StartTime.Format("2006-01-02 15:04"),
			s.EndTime.Format("2006-01-02 15:04"),
			durationStr(duration),
			xp,
		)
		fmt.Println(Green + line + Reset)
	}
}

func (app *application) printFunSessionsTable(sessions []*models.FunSession) {
	line := strings.Repeat("-", 93)
	fmt.Println(Green + line + Reset)
	fmt.Printf(Green+"|  %-3s |  %-15s |  %-16s |  %-16s |  %-11s |  %-7s |\n"+Reset,
		"ID", "App", "Start Time", "End Time", "Duration", "XP")
	fmt.Println(Green + line + Reset)

	for _, s := range sessions {
		duration := s.EndTime.Sub(s.StartTime)
		xp := app.calculateSpentXP(duration)
		fmt.Printf(
			pipe+"  %-3d "+pipe+"  "+Yellow+"%-15s"+Reset+" "+pipe+"  %-16s "+pipe+"  %-16s "+pipe+"  "+Yellow+"%-11s"+Reset+" "+pipe+"  "+Red+"%-7d"+Reset+" "+pipe+"\n",
			s.ID,
			s.App,
			s.StartTime.Format("2006-01-02 15:04"),
			s.EndTime.Format("2006-01-02 15:04"),
			durationStr(duration),
			xp,
		)
		fmt.Println(Green + line + Reset)
	}
}

func (app *application) printWorkTagSummaries(summaries []*models.WorkTagSummary) {
	line := strings.Repeat("-", 60)
	fmt.Println(Green + line + Reset)
	fmt.Printf(Green+"|  %-20s |  %-20s |  %-7s |\n"+Reset, "Tag", "Total Duration", "XP")
	fmt.Println(Green + line + Reset)

	for _, s := range summaries {
		xp := app.calculateEarnedXP(s.Duration)
		fmt.Printf(
			pipe+"  "+Yellow+"%-20s"+Reset+" "+pipe+"  "+Yellow+"%-20s"+Reset+" "+pipe+"  "+Red+"%-7d"+Reset+" "+pipe+"\n",
			s.Tag,
			durationStr(s.Duration),
			xp,
		)
		fmt.Println(Green + line + Reset)
	}
}

func (app *application) printFunAppSummaries(summaries []*models.FunAppSummary) {
	line := strings.Repeat("-", 60)
	fmt.Println(Green + line + Reset)
	fmt.Printf(Green+"|  %-20s |  %-20s |  %-7s |\n"+Reset, "App", "Total Duration", "XP")
	fmt.Println(Green + line + Reset)

	for _, s := range summaries {
		xp := app.calculateEarnedXP(s.Duration)
		fmt.Printf(
			pipe+"  "+Yellow+"%-20s"+Reset+" "+pipe+"  "+Yellow+"%-20s"+Reset+" "+pipe+"  "+Red+"%-7d"+Reset+" "+pipe+"\n",
			s.App,
			durationStr(s.Duration),
			xp,
		)
		fmt.Println(Green + line + Reset)
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
