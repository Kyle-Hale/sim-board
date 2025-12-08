package main

import "database/sql"

type App struct {
	DB *sql.DB
}

type Track struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Car       string `json:"car"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
}

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
}

type LapTime struct {
	ID               int    `json:"id"`
	TrackID          int    `json:"track_id"`
	UserID           int    `json:"user_id"`
	Username         string `json:"username"`
	Time             string `json:"time"`
	Car              string `json:"car"`
	ABS              bool   `json:"abs"`
	AutoTransmission bool   `json:"auto_transmission"`
	TractionControl  bool   `json:"traction_control"`
	CreatedAt        string `json:"created_at"`
}

type LeaderboardEntry struct {
	UserID           int    `json:"user_id"`
	Username         string `json:"username"`
	Time             string `json:"time"`
	Car              string `json:"car"`
	Position         int    `json:"position"`
	ABS              bool   `json:"abs"`
	AutoTransmission bool   `json:"auto_transmission"`
	TractionControl  bool   `json:"traction_control"`
	CreatedAt        string `json:"created_at"`
}
