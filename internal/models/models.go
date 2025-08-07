package models

import (
	"database/sql"
	"errors"
	"time"
)

var ErrRecordNotFound = errors.New("record not found")

type Models struct {
	WorkSessions interface {
		Insert(session *WorkSession) (int64, error)
		GetByID(id int64) (*WorkSession, error)
		FinishSession(id int64) error
		GetLastN(count int) ([]*WorkSession, error)
		TotalWorkDuration() (time.Duration, error)
		GetWorkTimeByTag(count int) ([]*WorkTagSummary, error)
		GetWeeklyWorkStats() (map[time.Weekday]time.Duration, error)
	}
	State interface {
		Get() (*UserState, error)
		Update(state *UserState) error
	}
	FunSessions interface {
		Insert(session *FunSession) error
		GetLastN(count int) ([]*FunSession, error)
		GetFunTimeByApp(count int) ([]*FunAppSummary, error)
	}

	LockedApps interface {
		GetAll() ([]*LockedApp, error)
	}
}

func NewModels(db *sql.DB) *Models {
	return &Models{
		WorkSessions: WorkSessionModel{DB: db},
		State:        UserStateModel{DB: db},
		FunSessions:  FunSessionModel{DB: db},
		LockedApps:   LockedAppModel{DB: db},
	}
}
