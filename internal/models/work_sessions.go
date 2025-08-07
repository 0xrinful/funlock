package models

import (
	"database/sql"
	"errors"
	"time"
)

type WorkSession struct {
	ID        int64
	Tag       string
	StartTime time.Time
	EndTime   time.Time
}

type WorkSessionModel struct {
	DB *sql.DB
}

func (m WorkSessionModel) Insert(session *WorkSession) (int64, error) {
	query := `
		INSERT INTO work_sessions (tag, start_time)
		VALUES (?, ?)`

	result, err := m.DB.Exec(query, session.Tag, session.StartTime)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m WorkSessionModel) GetByID(id int64) (*WorkSession, error) {
	query := `SELECT id, tag, start_time, end_time FROM work_sessions WHERE id = ?`
	var session WorkSession
	err := m.DB.QueryRow(query, id).
		Scan(&session.ID, &session.Tag, &session.StartTime, &session.EndTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &session, nil
}

func (m WorkSessionModel) FinishSession(id int64) error {
	query := `
		UPDATE work_sessions
		SET end_time = CURRENT_TIMESTAMP
		WHERE id = ?`
	_, err := m.DB.Exec(query, id)
	return err
}

func (m WorkSessionModel) GetLastN(count int) ([]*WorkSession, error) {
	query := `
		SELECT id, tag, start_time, end_time FROM work_sessions
		WHERE end_time != '0001-01-01 00:00:00'
		ORDER BY end_time DESC
		LIMIT ?`
	rows, err := m.DB.Query(query, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessions := make([]*WorkSession, 0, count)
	for rows.Next() {
		var session WorkSession
		err = rows.Scan(&session.ID, &session.Tag, &session.StartTime, &session.EndTime)
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

func (m WorkSessionModel) TotalWorkDuration() (time.Duration, error) {
	var totalSeconds *int64

	query := `
		SELECT SUM(strftime('%s', end_time) - strftime('%s', start_time))
		FROM work_sessions
		WHERE end_time != '0001-01-01 00:00:00'`

	err := m.DB.QueryRow(query).Scan(&totalSeconds)
	if err != nil {
		return 0, err
	}

	if totalSeconds == nil {
		return 0, nil
	}

	return time.Duration(*totalSeconds) * time.Second, nil
}

type WorkTagSummary struct {
	Tag      string
	Duration time.Duration
}

func (m WorkSessionModel) GetWorkTimeByTag(count int) ([]*WorkTagSummary, error) {
	query := `
		SELECT tag, SUM(strftime('%s', end_time) - strftime('%s', start_time)) AS total_seconds
		FROM work_sessions
		WHERE end_time != '0001-01-01 00:00:00'
		GROUP BY tag
		ORDER BY total_seconds DESC
		LIMIT ?`
	rows, err := m.DB.Query(query, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	summaries := make([]*WorkTagSummary, 0, count)

	for rows.Next() {
		var summary WorkTagSummary
		var totalSeconds int64
		if err := rows.Scan(&summary.Tag, &totalSeconds); err != nil {
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

func (m WorkSessionModel) GetWeeklyWorkStats() (map[time.Weekday]time.Duration, error) {
	query := `
		SELECT start_time, end_time FROM work_sessions
		WHERE end_time != '0001-01-01 00:00:00' AND start_time >= ?`

	now := time.Now()
	startOfWeek := time.Date(
		now.Year(),
		now.Month(),
		now.Day()-int(now.Weekday()),
		0, 0, 0, 0,
		now.Location(),
	)

	rows, err := m.DB.Query(query, startOfWeek.Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[time.Weekday]time.Duration)
	for rows.Next() {
		var startTime, endTime time.Time
		if err := rows.Scan(&startTime, &endTime); err != nil {
			return nil, err
		}
		duration := endTime.Sub(startTime)
		day := startTime.Weekday()
		result[day] += duration
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
