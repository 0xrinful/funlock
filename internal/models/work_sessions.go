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
