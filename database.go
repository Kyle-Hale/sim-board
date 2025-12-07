package main

import (
	"log"
	"strings"
)

func (app *App) initDB() {
	createTracksTable := `
	CREATE TABLE IF NOT EXISTS tracks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		car TEXT,
		is_active INTEGER DEFAULT 0,
		image_path TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	createLapTimesTable := `
	CREATE TABLE IF NOT EXISTS lap_times (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		track_id INTEGER NOT NULL,
		user_id INTEGER,
		username TEXT,
		time TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (track_id) REFERENCES tracks(id)
	);`

	if _, err := app.DB.Exec(createTracksTable); err != nil {
		log.Fatal("Failed to create tracks table:", err)
	}

	if _, err := app.DB.Exec(createUsersTable); err != nil {
		log.Fatal("Failed to create users table:", err)
	}

	if _, err := app.DB.Exec(createLapTimesTable); err != nil {
		log.Fatal("Failed to create lap_times table:", err)
	}

	app.DB.Exec("ALTER TABLE lap_times ADD COLUMN user_id INTEGER")
	app.DB.Exec("ALTER TABLE lap_times ADD COLUMN updated_at DATETIME")

	var usernameColExists int
	app.DB.QueryRow("SELECT COUNT(*) FROM pragma_table_info('lap_times') WHERE name='username'").Scan(&usernameColExists)
	if usernameColExists == 0 {
		app.DB.Exec("ALTER TABLE lap_times ADD COLUMN username TEXT")
	}

	var usernameCount int
	app.DB.QueryRow("SELECT COUNT(*) FROM pragma_table_info('lap_times') WHERE name='username'").Scan(&usernameCount)
	if usernameCount > 0 {
		rows, err := app.DB.Query("SELECT DISTINCT username FROM lap_times WHERE username IS NOT NULL AND username != '' AND (user_id IS NULL OR user_id = 0)")
		if err == nil {
			for rows.Next() {
				var username string
				if err := rows.Scan(&username); err != nil {
					continue
				}
				// Enforce uppercase for consistency
				upperUsername := strings.ToUpper(username)
				var userID int
				err := app.DB.QueryRow("SELECT id FROM users WHERE username = ?", upperUsername).Scan(&userID)
				if err != nil {
					result, err := app.DB.Exec("INSERT INTO users (username) VALUES (?)", upperUsername)
					if err != nil {
						continue
					}
					newID, err := result.LastInsertId()
					if err != nil {
						continue
					}
					userID = int(newID)
				}
				if userID > 0 {
					// Update both the old and new username format in lap_times
					app.DB.Exec("UPDATE lap_times SET user_id = ?, username = ? WHERE (username = ? OR username = ?) AND (user_id IS NULL OR user_id = 0)", userID, upperUsername, username, upperUsername)
				}
			}
			rows.Close()
		}
	}

	app.DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_track_user ON lap_times(track_id, user_id) WHERE user_id IS NOT NULL")

	createSettingsTable := `
	CREATE TABLE IF NOT EXISTS settings (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL
	);`

	if _, err := app.DB.Exec(createSettingsTable); err != nil {
		log.Fatal("Failed to create settings table:", err)
	}

	app.DB.Exec("ALTER TABLE tracks ADD COLUMN car TEXT")

	// Add assist columns to lap_times table
	var absColExists int
	app.DB.QueryRow("SELECT COUNT(*) FROM pragma_table_info('lap_times') WHERE name='abs'").Scan(&absColExists)
	if absColExists == 0 {
		app.DB.Exec("ALTER TABLE lap_times ADD COLUMN abs INTEGER DEFAULT 0")
	}

	var autoColExists int
	app.DB.QueryRow("SELECT COUNT(*) FROM pragma_table_info('lap_times') WHERE name='auto_transmission'").Scan(&autoColExists)
	if autoColExists == 0 {
		// Migrate from manual_transmission to auto_transmission if needed
		var mtColExists int
		app.DB.QueryRow("SELECT COUNT(*) FROM pragma_table_info('lap_times') WHERE name='manual_transmission'").Scan(&mtColExists)
		if mtColExists > 0 {
			// Rename column: manual = 0 means auto = 1, manual = 1 means auto = 0
			app.DB.Exec("ALTER TABLE lap_times ADD COLUMN auto_transmission INTEGER DEFAULT 0")
			app.DB.Exec("UPDATE lap_times SET auto_transmission = CASE WHEN manual_transmission = 0 THEN 1 ELSE 0 END")
		} else {
			app.DB.Exec("ALTER TABLE lap_times ADD COLUMN auto_transmission INTEGER DEFAULT 0")
		}
	}

	var tcColExists int
	app.DB.QueryRow("SELECT COUNT(*) FROM pragma_table_info('lap_times') WHERE name='traction_control'").Scan(&tcColExists)
	if tcColExists == 0 {
		app.DB.Exec("ALTER TABLE lap_times ADD COLUMN traction_control INTEGER DEFAULT 0")
	}

	var count int
	app.DB.QueryRow("SELECT COUNT(*) FROM tracks").Scan(&count)
	if count == 0 {
		app.DB.Exec("INSERT INTO tracks (name, is_active) VALUES (?, ?)", "Default Track", 1)
	}

	var settingCount int
	app.DB.QueryRow("SELECT COUNT(*) FROM settings WHERE key = 'show_admin_button'").Scan(&settingCount)
	if settingCount == 0 {
		app.DB.Exec("INSERT INTO settings (key, value) VALUES ('show_admin_button', 'true')")
	}

	var assistsSettingCount int
	app.DB.QueryRow("SELECT COUNT(*) FROM settings WHERE key = 'show_assists_leaderboard'").Scan(&assistsSettingCount)
	if assistsSettingCount == 0 {
		// Migrate from old setting name if it exists
		var oldSettingCount int
		app.DB.QueryRow("SELECT COUNT(*) FROM settings WHERE key = 'show_assists_admin'").Scan(&oldSettingCount)
		if oldSettingCount > 0 {
			var oldValue string
			app.DB.QueryRow("SELECT value FROM settings WHERE key = 'show_assists_admin'").Scan(&oldValue)
			app.DB.Exec("INSERT INTO settings (key, value) VALUES ('show_assists_leaderboard', ?)", oldValue)
		} else {
			app.DB.Exec("INSERT INTO settings (key, value) VALUES ('show_assists_leaderboard', 'true')")
		}
	}

	var titleSettingCount int
	app.DB.QueryRow("SELECT COUNT(*) FROM settings WHERE key = 'leaderboard_title'").Scan(&titleSettingCount)
	if titleSettingCount == 0 {
		app.DB.Exec("INSERT INTO settings (key, value) VALUES ('leaderboard_title', 'Sim Racing Leaderboard')")
	}

	// Migrate all existing usernames to uppercase for consistency
	rows, err := app.DB.Query("SELECT id, username FROM users")
	if err == nil {
		for rows.Next() {
			var userID int
			var username string
			if err := rows.Scan(&userID, &username); err != nil {
				continue
			}
			upperUsername := strings.ToUpper(username)
			if upperUsername != username {
				app.DB.Exec("UPDATE users SET username = ? WHERE id = ?", upperUsername, userID)
				// Also update username in lap_times if it exists
				app.DB.Exec("UPDATE lap_times SET username = ? WHERE username = ?", upperUsername, username)
			}
		}
		rows.Close()
	}
}
