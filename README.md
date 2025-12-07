# Sim-Board - Lap Time Leaderboard

A simple Go web application for tracking and displaying lap times on a racing simulator. Features a stylish leaderboard UI and an admin console for managing tracks and lap times.

## Features

- 🏁 **Leaderboard Display**: Beautiful, real-time leaderboard showing lap times for the active track
- ⚙️ **Admin Console**: Manage tracks, users, and lap times at `/admin`
- 🎯 **Active Track Selection**: Switch which track is displayed on the main leaderboard
- 🚗 **Sim Assists Tracking**: Track ABS, Auto Transmission, and Traction Control settings for each lap time
- 📊 **Customizable Display**: Toggle assist display on leaderboard, customize leaderboard title
- 👥 **User Management**: Add and manage users with automatic uppercase normalization
- 🐳 **Docker Support**: Easy deployment with Docker and Docker Compose
- 🌐 **Network Accessible**: Runs in a container accessible on your network
- 💚 **Health Checks**: Built-in health check endpoint for monitoring
- 🔒 **Production Ready**: Non-root user, resource limits, graceful shutdown, and connection pooling

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
2. **Track Management**:
   - Add tracks with optional car name
   - Set active track to display on leaderboard
   - Edit or delete tracks
3. **User Management**:
   - Add users (usernames are automatically converted to uppercase)
   - Edit or delete users
4. **Lap Time Management**:
   - Select a track and user
   - Enter lap time (format: `MM:SS.mmm`, e.g., `01:23.456`)
   - Optionally set sim assists (ABS, Auto Transmission, Traction Control)
   - Click "Add/Update Lap Time"
5. **Leaderboard Management**:
   - Customize leaderboard title
   - Toggle admin button visibility on leaderboard
   - Toggle assist display on leaderboard

### Leaderboard

The main page (`http://localhost:8869`) displays:
- Customizable title (default: "Sim Racing Leaderboard")
- The currently active track name and car
- A sorted leaderboard of all lap times for that track
- Top 3 positions are highlighted with special styling
- Sim assist indicators (ABS, Transmission, TC) - can be toggled in admin settings
- Assist status shown as: `ABS: ON/OFF | TRANS: AUTO/MANUAL | TC: ON/OFF`

## API Endpoints

### Tracks
- `GET /api/tracks` - Get all tracks
- `POST /api/tracks` - Create a new track
- `PUT /api/tracks/:id` - Update a track
- `DELETE /api/tracks/:id` - Delete a track
- `PUT /api/tracks/active` - Set the active track

### Users
- `GET /api/users` - Get all users
- `POST /api/users` - Create a new user (username auto-converted to uppercase)
- `PUT /api/users/:id` - Update a user (username auto-converted to uppercase)
- `DELETE /api/users/:id` - Delete a user

### Lap Times
- `GET /api/laptimes?track_id=X` - Get lap times for a track
- `POST /api/laptimes/add` - Add or update a lap time (includes sim assists)
- `PUT /api/laptimes/update` - Update a lap time
- `DELETE /api/laptimes/:id` - Delete a lap time

### Settings
- `GET /api/settings/admin-button` - Get admin button visibility setting
- `PUT /api/settings/admin-button` - Update admin button visibility
- `GET /api/settings/show-assists-leaderboard` - Get assist display setting
- `PUT /api/settings/show-assists-leaderboard` - Update assist display setting
- `GET /api/settings/leaderboard-title` - Get leaderboard title
- `PUT /api/settings/leaderboard-title` - Update leaderboard title

### Health
- `GET /health` - Health check endpoint (returns "healthy" or "unhealthy")

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

The SQLite database (`sim-board.db`) and uploaded images are stored in the container's data directory (`/app/data`). To persist data across container restarts, use a volume:

```bash
docker run -d -p 8869:8869 -v $(pwd)/data:/app/data --name sim-board sim-board
```

Or use the provided `docker-compose.yml` which includes volume mapping. The data directory location can be customized by setting the `DATA_DIR` environment variable before running docker-compose:

```bash
# Use a custom data directory
export DATA_DIR=/path/to/your/data
docker-compose up -d

# Or use the default ./data directory
docker-compose up -d
```

**Note**: The container runs as a non-root user (`appuser`) for security. The data directory is automatically created if it doesn't exist.

## Port Configuration

Change the port by setting the `PORT` environment variable:

```bash
docker run -d -p 9090:9090 -e PORT=9090 --name sim-board sim-board
```

## Sim Assists

Each lap time can include sim assist settings:
- **ABS**: Anti-lock Braking System (ON/OFF)
- **Auto Transmission**: Automatic or Manual transmission
- **Traction Control**: Traction Control (ON/OFF)

Assists are displayed on the leaderboard as text indicators (e.g., `ABS: ON | TRANS: AUTO | TC: OFF`). The display can be toggled in the admin console settings.

## Username Format

Usernames are automatically converted to uppercase for consistency. When creating or updating users, any case input will be normalized to uppercase (e.g., "john" becomes "JOHN").

## Health Check

The application includes a health check endpoint at `/health` that verifies database connectivity. This is used by Docker health checks and can be monitored externally.

## Production Features

- **Non-root user**: Container runs as `appuser` (UID 1000) for security
- **Resource limits**: CPU and memory limits configured in docker-compose
- **Graceful shutdown**: Handles SIGINT/SIGTERM signals properly
- **Connection pooling**: Database connection pool configured for optimal performance
- **HTTP timeouts**: Read, write, and idle timeouts configured
- **Health checks**: Built-in health monitoring

## Requirements

- Docker and Docker Compose (for containerized deployment)
- Go 1.21+ (for local development)
- SQLite (included in Docker image)

## License

MIT

