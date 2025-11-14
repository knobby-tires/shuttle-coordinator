package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

var apiKey = os.Getenv("FLIGHTAWARE_API_KEY")

// getFlightStatus fetches real-time flight data from FlightAware API
func getFlightStatus(flightNumber, flightType string, crewCount int) (Flight, error) {
	apiURL := fmt.Sprintf("https://aeroapi.flightaware.com/aeroapi/flights/%s", url.QueryEscape(flightNumber))

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return Flight{}, fmt.Errorf("Failed to create request")
	}

	// FlightAware uses x-apikey header for authentication
	req.Header.Set("x-apikey", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Flight{}, fmt.Errorf("Failed to connect to flight API")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Flight{}, fmt.Errorf("Failed to read API response")
	}

	// Parse JSON response from FlightAware
	var result struct {
		Flights []struct {
			Ident        string `json:"ident"`
			OperatorIata string `json:"operator_iata"`
			Operator     string `json:"operator"`
			Status       string `json:"status"`
			ScheduledIn  string `json:"scheduled_in"`
			EstimatedIn  string `json:"estimated_in"`
			ActualIn     string `json:"actual_in"`
		} `json:"flights"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return Flight{}, fmt.Errorf("Failed to parse API response: %s", err.Error())
	}

	if len(result.Flights) == 0 {
		return Flight{}, fmt.Errorf("Flight not found")
	}

	flightData := result.Flights[0]

	// Use most accurate time available (actual > estimated > scheduled)
	arrivalTimeStr := flightData.ScheduledIn
	if flightData.EstimatedIn != "" {
		arrivalTimeStr = flightData.EstimatedIn
	}
	if flightData.ActualIn != "" {
		arrivalTimeStr = flightData.ActualIn
	}

	if arrivalTimeStr == "" {
		return Flight{}, fmt.Errorf("Flight has no arrival time data available")
	}

	// Parse ISO 8601 time format
	scheduledTime, parseErr := time.Parse(time.RFC3339, arrivalTimeStr)
	if parseErr != nil {
		return Flight{}, fmt.Errorf("Failed to parse time: %s", parseErr.Error())
	}

	// Convert to Mountain Time (MST/MDT)
	mountainTime, _ := time.LoadLocation("America/Denver")
	scheduledMT := scheduledTime.In(mountainTime)

	// Calculate delay by comparing scheduled vs estimated
	var delay int
	if flightData.ScheduledIn != "" && flightData.EstimatedIn != "" {
		scheduledT, _ := time.Parse(time.RFC3339, flightData.ScheduledIn)
		estimatedT, _ := time.Parse(time.RFC3339, flightData.EstimatedIn)
		delay = int(estimatedT.Sub(scheduledT).Minutes())
	}

	expectedMT := scheduledMT

	// Prefer full airline name over IATA code
	airlineName := flightData.Operator
	if airlineName == "" {
		airlineName = flightData.OperatorIata
	}

	flight := Flight{
		FlightNumber:     flightData.Ident,
		Airline:          airlineName,
		Status:           flightData.Status,
		ScheduledArrival: scheduledMT.Format("3:04 PM"),
		ExpectedArrival:  expectedMT.Format("3:04 PM"),
		Delay:            delay,
		IsDelayed:        delay > 0,
		Type:             flightType,
		CrewCount:        crewCount,
		Note:             "",
		SortTime:         expectedMT,
	}

	return flight, nil
}
