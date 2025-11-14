package main

import "time"

// Flight represents a single flight with shuttle coordination details
type Flight struct {
	ID               int       // Unique identifier for this flight
	FlightNumber     string    // Flight number (e.g., "AA100")
	Airline          string    // Airline name
	Status           string    // Flight status (scheduled, active, landed, etc.)
	ScheduledArrival string    // Original scheduled arrival time (formatted)
	ExpectedArrival  string    // Expected arrival time accounting for delays (formatted)
	Delay            int       // Delay in minutes
	IsDelayed        bool      // Whether the flight is delayed
	Type             string    // "pickup", "dropoff", or "both"
	CrewCount        int       // Number of crew members to transport
	Note             string    // Optional note for this flight
	SortTime         time.Time // Used for sorting flights chronologically
}

// PageData is the data passed to the HTML template
type PageData struct {
	Flights []Flight
	Error   string
	IsDemo  bool // Whether the current user is a demo account
}
