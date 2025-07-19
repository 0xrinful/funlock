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

func (app *application) calculateXP(duration time.Duration) int64 {
	return int64(duration.Seconds() * app.config.XPFactor)
}

func (app *application) printWorkSessionsTable(sessions []*models.WorkSession) {
	line := strings.Repeat("-", 93)
	fmt.Println(line)
	fmt.Printf("|  %-3s |  %-15s |  %-16s |  %-16s |  %-11s |  %-7s |\n",
		"ID", "Tag", "Start Time", "End Time", "Duration", "XP")
	fmt.Println(line)

	for _, s := range sessions {
		duration := s.EndTime.Sub(s.StartTime)
		fmt.Printf("|  %-3d |  %-15s |  %-16s |  %-16s |  %-11s |  %-7d |\n",
			s.ID,
			s.Tag,
			s.StartTime.Format("2006-01-02 15:04"),
			s.EndTime.Format("2006-01-02 15:04"),
			durationStr(duration),
			app.calculateXP(duration),
		)
		fmt.Println(line)
	}
}

func (app *application) printFunSessionsTable(sessions []*models.FunSession) {
	line := strings.Repeat("-", 93)
	fmt.Println(line)

	fmt.Printf("|  %-3s |  %-15s |  %-16s |  %-16s |  %-11s |  %-7s |\n",
		"ID", "App", "Start Time", "End Time", "Duration", "XP")
	fmt.Println(line)

	for _, s := range sessions {
		duration := s.EndTime.Sub(s.StartTime)
		fmt.Printf("|  %-3d |  %-15s |  %-16s |  %-16s |  %-11s |  %-7d |\n",
			s.ID,
			s.App,
			s.StartTime.Format("2006-01-02 15:04"),
			s.EndTime.Format("2006-01-02 15:04"),
			durationStr(duration),
			app.calculateXP(duration),
		)
		fmt.Println(line)
	}
}
