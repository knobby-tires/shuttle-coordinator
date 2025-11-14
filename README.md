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

## Project Structure
```
shuttle-coordinator/
├── main.go           # Server and HTTP handlers
├── auth.go           # Authentication & sessions
├── flightaware.go    # FlightAware API client
├── models.go         # Data structures
├── template.go       # HTML templates
├── tests/
│   ├── auth_test.go      # Authentication tests
│   ├── session_test.go   # Session management tests
│   └── password_test.go  # Password hashing tests
├── screenshots/          # UI examples
├── go.mod            # Go dependencies
└── README.md
```

## Key Features Explained

### Authentication & Security
- Bcrypt password hashing with configurable cost factor
- HTTP-only session cookies prevent XSS attacks
- SameSite cookie policy prevents CSRF
- Middleware-protected routes
- Role-based access (valet, desk, demo)

### Flight Tracking
- Real-time data from FlightAware AeroAPI
- Automatic timezone conversion to Mountain Time
- Delay calculation and visual indicators
- Scheduled vs. expected arrival times
- Flight status monitoring (scheduled, active, landed)

### Demo Mode
- Simulated flight data for demonstrations
- No API calls made (prevents costs)
- Showcases full functionality without live data

### Testing
Comprehensive test suite covering:
- Password hashing and verification
- Session lifecycle management
- Authentication flow
- Middleware protection
- Concurrent session access

Run tests with: `go test -v`

## Design Decisions

- **In-memory storage**: Simplifies deployment; production would use PostgreSQL
- **Minimal JavaScript**: Ensures compatibility with older iPad devices
- **Session-based auth**: Simpler than JWT for this use case
- **Server-side rendering**: Better performance on older hardware

## Future Improvements

- Persistent database storage (PostgreSQL)
- Traffic integration for departure time calculations
- Email/SMS notifications for flight delays
- Multi-location support
- Historical flight data and analytics

## License

MIT
