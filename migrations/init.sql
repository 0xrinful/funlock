CREATE TABLE work_sessions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  tag TEXT,
  start_time DATETIME NOT NULL,
  end_time DATETIME NOT NULL DEFAULT '0001-01-01 00:00:00'
);

CREATE TABLE fun_sessions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  app TEXT NOT NULL,
  start_time DATETIME NOT NULL,
  end_time DATETIME NOT NULL
);

CREATE TABLE locked_apps (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL UNIQUE,
  required_xp INTEGER NOT NULL DEFAULT 0,
  permanently_locked BOOLEAN DEFAULT 0
);

CREATE TABLE user_state (
  id INTEGER PRIMARY KEY CHECK (id = 1),
  xp_balance INTEGER NOT NULL DEFAULT 0,
  current_session_id INTEGER REFERENCES work_sessions (id)
);

INSERT INTO
  user_state (id, xp_balance, current_session_id)
VALUES
  (1, 0, NULL);
