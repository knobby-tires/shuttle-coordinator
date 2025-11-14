# Flight Tracker

A real-time flight tracking and shuttle coordination system for managing airport transportation logistics. Tracks flight arrivals/departures, manages crew pickups, and coordinates shuttles with live flight delay information.

## Features

- Real-time flight tracking via FlightAware AeroAPI
- Secure authentication with bcrypt password hashing
- Role-based access control
- Automatic timezone conversion to Mountain Time
- Flight-specific notes for contextual information
- Pickup/dropoff management with crew counting
- Auto-refresh for live status updates
- Demo mode with simulated data

## Tech Stack

- Go (Golang)
- FlightAware AeroAPI
- bcrypt authentication
- HTML/CSS/JavaScript

## Installation

1. Clone the repository
```bash
git clone https://github.com/yourusername/flight-tracker.git
cd flight-tracker
```

2. Install dependencies
```bash
go mod init shuttletracker
go get golang.org/x/crypto/bcrypt
```

3. Set environment variables
```bash
export FLIGHTAWARE_API_KEY=your_api_key_here
export VALET_PASSWORD=your_secure_password
export DESK_PASSWORD=your_other_password
```

4. Run the application
```bash
go run *.go
```

5. Access at `http://localhost:8080`

## Project Structure
```
ShuttleTracker/
├── main.go           # Server and HTTP handlers
├── auth.go           # Authentication & sessions
├── flightaware.go    # API client
├── models.go         # Data structures
├── template.go       # HTML templates
├── auth_test.go      # Authentication tests
├── session_test.go   # Session management tests
├── password_test.go  # Password hashing tests
└── go.mod            # Dependencies
```

## Testing

Run the test suite:
```bash
go test -v
```

## Authentication

The system uses bcrypt-hashed passwords with secure session cookies. Multiple user roles are supported for different operational needs.

## Design Decisions

- In-memory storage for simplicity (production would use PostgreSQL)
- Minimal JavaScript for compatibility with older devices
- Session-based authentication over JWT for reduced complexity
- Demo mode prevents API costs during demonstrations

## Security

- Bcrypt password hashing
- HTTP-only session cookies
- SameSite cookie policy
- Middleware-protected routes

## Future Improvements

- Persistent database storage
- Traffic integration for departure time calculations
- Notification system for delays
- Multi-location support
- Historical data analytics

## License

MIT
