package models

import (
	"database/sql"
	"time"
)

type FunSession struct {
	ID        int64
	App       string
	StartTime time.Time
	EndTime   time.Time
}

type FunSessionModel struct {
	DB *sql.DB
}

func (m FunSessionModel) Insert(session *FunSession) error {
	query := `
		INSERT INTO fun_sessions (app, start_time, end_time)
		VALUES (?, ?, ?)`

	_, err := m.DB.Exec(query, session.App, session.StartTime, session.EndTime)
	return err
}

func (m FunSessionModel) GetLastN(count int) ([]*FunSession, error) {
	query := `
		SELECT id, app, start_time, end_time FROM fun_sessions
		ORDER BY end_time DESC
		LIMIT ?`
	rows, err := m.DB.Query(query, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessions := make([]*FunSession, 0, count)
	for rows.Next() {
		var session FunSession
		err = rows.Scan(&session.ID, &session.App, &session.StartTime, &session.EndTime)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, &session)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sessions, nil
}

type FunAppSummary struct {
	App      string
	Duration time.Duration
}

func (m FunSessionModel) GetFunTimeByApp(count int) ([]*FunAppSummary, error) {
	query := `
		SELECT app, SUM(strftime('%s', end_time) - strftime('%s', start_time)) AS total_seconds
		FROM fun_sessions
		GROUP BY app
		ORDER BY total_seconds DESC
		LIMIT ?`
	rows, err := m.DB.Query(query, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	summaries := make([]*FunAppSummary, 0, count)

	for rows.Next() {
		var summary FunAppSummary
		var totalSeconds int64
		if err := rows.Scan(&summary.App, &totalSeconds); err != nil {
			return nil, err
		}
		summary.Duration = time.Duration(totalSeconds) * time.Second
		summaries = append(summaries, &summary)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return summaries, nil
}
