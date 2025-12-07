package main

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"net/http"
	"strings"
)

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func (app *App) handleHealth(w http.ResponseWriter, r *http.Request) {
	// Simple health check - verify database connection
	if err := app.DB.Ping(); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("unhealthy"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("healthy"))
}

func (app *App) handleLeaderboard(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("leaderboard").Parse(leaderboardHTML))

	var activeTrack Track
	var car sql.NullString
	err := app.DB.QueryRow(`
		SELECT id, name, car, is_active, created_at 
		FROM tracks 
		WHERE is_active = 1 
		LIMIT 1
	`).Scan(&activeTrack.ID, &activeTrack.Name, &car, &activeTrack.IsActive, &activeTrack.CreatedAt)

	if err != nil {
		activeTrack = Track{ID: 0, Name: "No Active Track"}
	} else {
		if car.Valid {
			activeTrack.Car = car.String
		}
	}

	rows, err := app.DB.Query(`
		SELECT u.id, u.username, lt.time, 
		       COALESCE(lt.abs, 0) as abs,
		       COALESCE(lt.auto_transmission, 0) as auto_transmission,
		       COALESCE(lt.traction_control, 0) as traction_control,
		       lt.created_at 
		FROM lap_times lt
		JOIN users u ON lt.user_id = u.id
		WHERE lt.track_id = ? 
		ORDER BY 
			CAST(SUBSTR(lt.time, 1, 2) AS INTEGER) ASC,
			CAST(SUBSTR(lt.time, 4, 2) AS INTEGER) ASC,
			CAST(SUBSTR(lt.time, 7) AS REAL) ASC
		LIMIT 100
	`, activeTrack.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var entries []LeaderboardEntry
	position := 1
	for rows.Next() {
		var entry LeaderboardEntry
		var abs, mt, tc int
		rows.Scan(&entry.UserID, &entry.Username, &entry.Time, &abs, &mt, &tc, &entry.CreatedAt)
		entry.ABS = abs == 1
		entry.AutoTransmission = mt == 1
		entry.TractionControl = tc == 1
		entry.Position = position
		entries = append(entries, entry)
		position++
	}

	var showAdminButton bool
	var settingValue string
	err = app.DB.QueryRow("SELECT value FROM settings WHERE key = 'show_admin_button'").Scan(&settingValue)
	if err != nil {
		showAdminButton = true
	} else {
		showAdminButton = settingValue == "true"
	}

	var leaderboardTitle string
	err = app.DB.QueryRow("SELECT value FROM settings WHERE key = 'leaderboard_title'").Scan(&leaderboardTitle)
	if err != nil {
		leaderboardTitle = "Sim Racing Leaderboard"
	}

	var showAssistsLeaderboard bool
	var assistsSettingValue string
	err = app.DB.QueryRow("SELECT value FROM settings WHERE key = 'show_assists_leaderboard'").Scan(&assistsSettingValue)
	if err != nil {
		showAssistsLeaderboard = true
	} else {
		showAssistsLeaderboard = assistsSettingValue == "true"
	}

	data := struct {
		Track                  Track
		Entries                []LeaderboardEntry
		ShowAdminButton        bool
		ShowAssistsLeaderboard bool
		Title                  string
	}{
		Track:                  activeTrack,
		Entries:                entries,
		ShowAdminButton:        showAdminButton,
		ShowAssistsLeaderboard: showAssistsLeaderboard,
		Title:                  leaderboardTitle,
	}

	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, data)
}

func (app *App) handleAdmin(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("admin").Parse(adminHTML))
	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, nil)
}

func (app *App) handleTracks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		rows, err := app.DB.Query("SELECT id, name, car, is_active, created_at FROM tracks ORDER BY created_at DESC")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var tracks []Track
		for rows.Next() {
			var track Track
			var isActive int
			var car sql.NullString
			rows.Scan(&track.ID, &track.Name, &car, &isActive, &track.CreatedAt)
			track.IsActive = isActive == 1
			if car.Valid {
				track.Car = car.String
			}
			tracks = append(tracks, track)
		}

		json.NewEncoder(w).Encode(tracks)
		return
	}

	if r.Method == "POST" {
		var track Track
		json.NewDecoder(r.Body).Decode(&track)

		result, err := app.DB.Exec("INSERT INTO tracks (name, car, is_active) VALUES (?, ?, ?)", track.Name, track.Car, 0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		id, _ := result.LastInsertId()
		track.ID = int(id)
		json.NewEncoder(w).Encode(track)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (app *App) handleTrackByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := r.URL.Path
	trackID := strings.TrimPrefix(path, "/api/tracks/")
	if trackID == "" || trackID == "active" {
		http.Error(w, "Invalid track ID", http.StatusBadRequest)
		return
	}

	if r.Method == "PUT" {
		var track Track
		json.NewDecoder(r.Body).Decode(&track)

		result, err := app.DB.Exec("UPDATE tracks SET name = ?, car = ? WHERE id = ?", track.Name, track.Car, trackID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, "Track not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
		return
	}

	if r.Method == "DELETE" {
		app.DB.Exec("DELETE FROM lap_times WHERE track_id = ?", trackID)

		result, err := app.DB.Exec("DELETE FROM tracks WHERE id = ?", trackID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, "Track not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (app *App) handleActiveTrack(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "PUT" {
		var data struct {
			TrackID int `json:"track_id"`
		}
		json.NewDecoder(r.Body).Decode(&data)

		app.DB.Exec("UPDATE tracks SET is_active = 0")
		_, err := app.DB.Exec("UPDATE tracks SET is_active = 1 WHERE id = ?", data.TrackID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (app *App) handleUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		rows, err := app.DB.Query("SELECT id, username, created_at FROM users ORDER BY username ASC")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []User
		for rows.Next() {
			var user User
			rows.Scan(&user.ID, &user.Username, &user.CreatedAt)
			users = append(users, user)
		}

		json.NewEncoder(w).Encode(users)
		return
	}

	if r.Method == "POST" {
		var user User
		json.NewDecoder(r.Body).Decode(&user)

		// Enforce uppercase for consistency
		user.Username = strings.ToUpper(strings.TrimSpace(user.Username))

		result, err := app.DB.Exec("INSERT INTO users (username) VALUES (?)", user.Username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		id, _ := result.LastInsertId()
		user.ID = int(id)
		json.NewEncoder(w).Encode(user)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (app *App) handleUserByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := r.URL.Path
	userID := strings.TrimPrefix(path, "/api/users/")
	if userID == "" {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if r.Method == "PUT" {
		var user User
		json.NewDecoder(r.Body).Decode(&user)

		// Enforce uppercase for consistency
		user.Username = strings.ToUpper(strings.TrimSpace(user.Username))

		result, err := app.DB.Exec("UPDATE users SET username = ? WHERE id = ?", user.Username, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
		return
	}

	if r.Method == "DELETE" {
		app.DB.Exec("DELETE FROM lap_times WHERE user_id = ?", userID)

		result, err := app.DB.Exec("DELETE FROM users WHERE id = ?", userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (app *App) handleLapTimes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	trackID := r.URL.Query().Get("track_id")
	if trackID == "" {
		http.Error(w, "track_id required", http.StatusBadRequest)
		return
	}

	rows, err := app.DB.Query(`
		SELECT lt.id, lt.track_id, lt.user_id, u.username, lt.time,
		       COALESCE(lt.abs, 0) as abs,
		       COALESCE(lt.auto_transmission, 0) as auto_transmission,
		       COALESCE(lt.traction_control, 0) as traction_control,
		       lt.created_at 
		FROM lap_times lt
		JOIN users u ON lt.user_id = u.id
		WHERE lt.track_id = ? 
		ORDER BY 
			CAST(SUBSTR(lt.time, 1, 2) AS INTEGER) ASC,
			CAST(SUBSTR(lt.time, 4, 2) AS INTEGER) ASC,
			CAST(SUBSTR(lt.time, 7) AS REAL) ASC
	`, trackID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var lapTimes []LapTime
	for rows.Next() {
		var lt LapTime
		var abs, mt, tc int
		rows.Scan(&lt.ID, &lt.TrackID, &lt.UserID, &lt.Username, &lt.Time, &abs, &mt, &tc, &lt.CreatedAt)
		lt.ABS = abs == 1
		lt.AutoTransmission = mt == 1
		lt.TractionControl = tc == 1
		lapTimes = append(lapTimes, lt)
	}

	json.NewEncoder(w).Encode(lapTimes)
}

func (app *App) handleAddLapTime(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		TrackID          int    `json:"track_id"`
		UserID           int    `json:"user_id"`
		Time             string `json:"time"`
		ABS              bool   `json:"abs"`
		AutoTransmission bool   `json:"auto_transmission"`
		TractionControl  bool   `json:"traction_control"`
	}
	json.NewDecoder(r.Body).Decode(&data)

	var existingID int
	err := app.DB.QueryRow(
		"SELECT id FROM lap_times WHERE track_id = ? AND user_id = ?",
		data.TrackID, data.UserID,
	).Scan(&existingID)

	if err == nil {
		_, err = app.DB.Exec(
			"UPDATE lap_times SET time = ?, abs = ?, auto_transmission = ?, traction_control = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
			data.Time, boolToInt(data.ABS), boolToInt(data.AutoTransmission), boolToInt(data.TractionControl), existingID,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
		return
	}

	_, err = app.DB.Exec(
		"INSERT INTO lap_times (track_id, user_id, time, abs, auto_transmission, traction_control, updated_at, username) VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, NULL)",
		data.TrackID, data.UserID, data.Time, boolToInt(data.ABS), boolToInt(data.AutoTransmission), boolToInt(data.TractionControl),
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (app *App) handleUpdateLapTime(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "PUT" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		TrackID          int    `json:"track_id"`
		UserID           int    `json:"user_id"`
		Time             string `json:"time"`
		ABS              bool   `json:"abs"`
		AutoTransmission bool   `json:"auto_transmission"`
		TractionControl  bool   `json:"traction_control"`
	}
	json.NewDecoder(r.Body).Decode(&data)

	result, err := app.DB.Exec(
		"UPDATE lap_times SET time = ?, abs = ?, auto_transmission = ?, traction_control = ?, updated_at = CURRENT_TIMESTAMP WHERE track_id = ? AND user_id = ?",
		data.Time, boolToInt(data.ABS), boolToInt(data.AutoTransmission), boolToInt(data.TractionControl), data.TrackID, data.UserID,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Lap time not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (app *App) handleLapTimeByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := r.URL.Path
	lapTimeID := strings.TrimPrefix(path, "/api/laptimes/")
	if lapTimeID == "" || lapTimeID == "add" {
		http.Error(w, "Invalid lap time ID", http.StatusBadRequest)
		return
	}

	if r.Method == "DELETE" {
		result, err := app.DB.Exec("DELETE FROM lap_times WHERE id = ?", lapTimeID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, "Lap time not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (app *App) handleAdminButtonSetting(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		var value string
		err := app.DB.QueryRow("SELECT value FROM settings WHERE key = 'show_admin_button'").Scan(&value)
		if err != nil {
			value = "true"
		}
		json.NewEncoder(w).Encode(map[string]string{"value": value})
		return
	}

	if r.Method == "PUT" {
		var data struct {
			Value string `json:"value"`
		}
		json.NewDecoder(r.Body).Decode(&data)

		_, err := app.DB.Exec("INSERT OR REPLACE INTO settings (key, value) VALUES ('show_admin_button', ?)", data.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (app *App) handleLeaderboardTitleSetting(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		var value string
		err := app.DB.QueryRow("SELECT value FROM settings WHERE key = 'leaderboard_title'").Scan(&value)
		if err != nil {
			value = "Sim Racing Leaderboard"
		}
		json.NewEncoder(w).Encode(map[string]string{"value": value})
		return
	}

	if r.Method == "PUT" {
		var data struct {
			Value string `json:"value"`
		}
		json.NewDecoder(r.Body).Decode(&data)

		_, err := app.DB.Exec("INSERT OR REPLACE INTO settings (key, value) VALUES ('leaderboard_title', ?)", data.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (app *App) handleShowAssistsLeaderboardSetting(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		var value string
		err := app.DB.QueryRow("SELECT value FROM settings WHERE key = 'show_assists_leaderboard'").Scan(&value)
		if err != nil {
			value = "true"
		}
		json.NewEncoder(w).Encode(map[string]string{"value": value})
		return
	}

	if r.Method == "PUT" {
		var data struct {
			Value string `json:"value"`
		}
		json.NewDecoder(r.Body).Decode(&data)

		_, err := app.DB.Exec("INSERT OR REPLACE INTO settings (key, value) VALUES ('show_assists_leaderboard', ?)", data.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
