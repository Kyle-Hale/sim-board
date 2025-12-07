# Sim-Board - Lap Time Leaderboard

A simple Go web application for tracking and displaying lap times on a racing simulator. Features a stylish leaderboard UI and an admin console for managing tracks and lap times.

## Features

- 🏁 **Leaderboard Display**: Beautiful, real-time leaderboard showing lap times for the active track
- ⚙️ **Admin Console**: Manage tracks and lap times at `/admin`
- 🎯 **Active Track Selection**: Switch which track is displayed on the main leaderboard
- 🐳 **Docker Support**: Easy deployment with Docker and Docker Compose
- 🌐 **Network Accessible**: Runs in a container accessible on your network

## Quick Start

### Using Docker Compose (Recommended)

1. Clone or navigate to the project directory:
```bash
cd sim-board
```

2. Build and run with Docker Compose:
```bash
docker-compose up -d
```

3. Access the application:
   - **Leaderboard**: http://localhost:8869
   - **Admin Console**: http://localhost:8869/admin

### Using Docker

1. Build the image:
```bash
docker build -t sim-board .
```

2. Run the container:
```bash
docker run -d -p 8869:8869 --name sim-board sim-board
```

### Local Development

1. Install Go dependencies:
```bash
go mod download
```

2. Run the application:
```bash
go run main.go
```

The application will start on port 8869 (or the port specified in the `PORT` environment variable).

## Usage

### Admin Console

1. Navigate to `http://localhost:8869/admin`
2. **Add Tracks**: Enter a track name and click "Add Track"
3. **Set Active Track**: Click "Set Active" on any track to display it on the main leaderboard
4. **Add Lap Times**: 
   - Select a track
   - Enter username and lap time (format: `MM:SS.mmm`, e.g., `01:23.456`)
   - Click "Add Lap Time"

### Leaderboard

The main page (`http://localhost:8869`) displays:
- The currently active track name
- A sorted leaderboard of all lap times for that track
- Top 3 positions are highlighted with special styling

## API Endpoints

- `GET /api/tracks` - Get all tracks
- `POST /api/tracks` - Create a new track
- `PUT /api/tracks/active` - Set the active track
- `GET /api/laptimes?track_id=X` - Get lap times for a track
- `POST /api/laptimes/add` - Add a new lap time

## Network Access

To make the application accessible on your network:

1. **Docker Compose**: The container is already configured with `network_mode: bridge`, making it accessible on your local network at `http://<your-ip>:8869`

2. **Docker Run**: Use:
```bash
docker run -d -p 0.0.0.0:8869:8869 --name sim-board sim-board
```

3. Find your IP address:
   - macOS/Linux: `ifconfig | grep "inet "`
   - Windows: `ipconfig`

Then access from other devices on your network using `http://<your-ip>:8869`

## Data Persistence

The SQLite database (`sim-board.db`) is stored in the container. To persist data across container restarts, use a volume:

```bash
docker run -d -p 8869:8869 -v $(pwd)/data:/root/data --name sim-board sim-board
```

Or use the provided `docker-compose.yml` which includes volume mapping. The data directory location can be customized by setting the `DATA_DIR` environment variable before running docker-compose:

```bash
# Use a custom data directory
export DATA_DIR=/path/to/your/data
docker-compose up -d

# Or use the default ./data directory
docker-compose up -d
```

## Port Configuration

Change the port by setting the `PORT` environment variable:

```bash
docker run -d -p 9090:9090 -e PORT=9090 --name sim-board sim-board
```

## Requirements

- Docker and Docker Compose (for containerized deployment)
- Go 1.21+ (for local development)
- SQLite (included in Docker image)

## License

MIT

