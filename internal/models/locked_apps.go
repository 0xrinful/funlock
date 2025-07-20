package models

import "database/sql"

type LockedApp struct {
	ID   int64
	Name string
}

type LockedAppModel struct {
	DB *sql.DB
}

func (m LockedAppModel) GetAll() ([]*LockedApp, error) {
	query := "SELECT id, name FROM locked_apps"

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	apps := []*LockedApp{}
	for rows.Next() {
		var app LockedApp
		err = rows.Scan(&app.ID, &app.Name)
		if err != nil {
			return nil, err
		}
		apps = append(apps, &app)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return apps, nil
}
