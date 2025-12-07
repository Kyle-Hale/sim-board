package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "./data"
	}
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		os.MkdirAll(dataDir, 0755)
	}

	db, err := sql.Open("sqlite3", dataDir+"/sim-board.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	// Configure database connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify database connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	app := &App{DB: db}
	app.initDB()

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/health", app.handleHealth)
	http.HandleFunc("/", app.handleLeaderboard)
	http.HandleFunc("/admin", app.handleAdmin)
	http.HandleFunc("/api/tracks/active", app.handleActiveTrack)
	http.HandleFunc("/api/tracks/", app.handleTrackByID)
	http.HandleFunc("/api/tracks", app.handleTracks)
	http.HandleFunc("/api/users/", app.handleUserByID)
	http.HandleFunc("/api/users", app.handleUsers)
	http.HandleFunc("/api/laptimes/add", app.handleAddLapTime)
	http.HandleFunc("/api/laptimes/update", app.handleUpdateLapTime)
	http.HandleFunc("/api/laptimes/", app.handleLapTimeByID)
	http.HandleFunc("/api/laptimes", app.handleLapTimes)
	http.HandleFunc("/api/settings/admin-button", app.handleAdminButtonSetting)
	http.HandleFunc("/api/settings/show-assists-leaderboard", app.handleShowAssistsLeaderboardSetting)
	http.HandleFunc("/api/settings/leaderboard-title", app.handleLeaderboardTitleSetting)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8869"
	}

	// Create HTTP server with timeouts
	srv := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", port)
		log.Printf("Leaderboard: http://localhost:%s", port)
		log.Printf("Admin: http://localhost:%s/admin", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
