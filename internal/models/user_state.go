package models

import "database/sql"

type UserState struct {
	XpBalance        int64
	CurrentSessionID *int64
}

func (s *UserState) IsSessionRunning() bool {
	return s.CurrentSessionID != nil
}

type UserStateModel struct {
	DB *sql.DB
}

func (m UserStateModel) Get() (*UserState, error) {
	query := `
		SELECT xp_balance, current_session_id
		FROM user_state
		where id = 1`
	var state UserState
	err := m.DB.QueryRow(query).
		Scan(&state.XpBalance, &state.CurrentSessionID)
	if err != nil {
		return nil, err
	}
	return &state, nil
}

func (m UserStateModel) Update(state *UserState) error {
	query := `
		UPDATE user_state
		SET current_session_id = ?,
		xp_balance = ?
		WHERE id = 1`
	_, err := m.DB.Exec(query, state.CurrentSessionID, state.XpBalance)
	return err
}
