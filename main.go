package main

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"time"
)

var flights = []Flight{}
var nextID = 1
var tmpl *template.Template

func main() {
	var err error
	tmpl, err = template.New("index").Parse(htmlTemplate)
	if err != nil {
		panic(err)
	}

	tmpl, err = tmpl.New("login").Parse(loginTemplate)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/", requireAuth(homeHandler))
	http.HandleFunc("/add", requireAuth(addFlightHandler))
	http.HandleFunc("/remove", requireAuth(removeFlightHandler))
	http.HandleFunc("/update-note", requireAuth(updateNoteHandler))
	http.HandleFunc("/logout", requireAuth(logoutHandler))

	fmt.Println("🛫 Jacob's Flight Tracker")
	fmt.Println("Server starting on http://localhost:8080")
	fmt.Println("\nDemo account: demo / demo123")
	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	user := getCurrentUser(r)
	isDemo := user != nil && user.Role == "demo"

	sortedFlights := make([]Flight, len(flights))
	copy(sortedFlights, flights)
	sort.Slice(sortedFlights, func(i, j int) bool {
		return sortedFlights[i].SortTime.Before(sortedFlights[j].SortTime)
	})

	tmpl.ExecuteTemplate(w, "index", PageData{
		Flights: sortedFlights,
		IsDemo:  isDemo,
	})
}

func addFlightHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user := getCurrentUser(r)
	isDemo := user != nil && user.Role == "demo"

	flightNumber := r.FormValue("flight_number")
	isPickup := r.FormValue("is_pickup") == "on"
	isDropoff := r.FormValue("is_dropoff") == "on"
	crewCount, _ := strconv.Atoi(r.FormValue("crew_count"))

	flightType := "pickup"
	if isPickup && isDropoff {
		flightType = "both"
	} else if isDropoff {
		flightType = "dropoff"
	}

	var flight Flight
	var err error

	if isDemo {
		flight = createDemoFlight(flightNumber, flightType, crewCount)
	} else {
		flight, err = getFlightStatus(flightNumber, flightType, crewCount)
		if err != nil {
			sortedFlights := make([]Flight, len(flights))
			copy(sortedFlights, flights)
			sort.Slice(sortedFlights, func(i, j int) bool {
				return sortedFlights[i].SortTime.Before(sortedFlights[j].SortTime)
			})
			tmpl.ExecuteTemplate(w, "index", PageData{
				Flights: sortedFlights,
				Error:   err.Error() + " (Note: Free API tier may not include all flights)",
				IsDemo:  isDemo,
			})
			return
		}
	}

	flight.ID = nextID
	nextID++
	flights = append(flights, flight)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func createDemoFlight(flightNumber, flightType string, crewCount int) Flight {
	now := time.Now()
	mountainTime, _ := time.LoadLocation("America/Denver")

	arrivalTime := now.Add(time.Hour * time.Duration(2+nextID%3)).In(mountainTime)

	delay := 0
	isDelayed := false
	if nextID%3 == 0 {
		delay = 15 + (nextID % 30)
		isDelayed = true
	}

	airlines := []string{"American Airlines", "Delta Air Lines", "United Airlines", "Southwest Airlines"}
	statuses := []string{"scheduled", "active"}

	return Flight{
		FlightNumber:     flightNumber,
		Airline:          airlines[nextID%len(airlines)],
		Status:           statuses[nextID%len(statuses)],
		ScheduledArrival: arrivalTime.Format("3:04 PM"),
		ExpectedArrival:  arrivalTime.Add(time.Duration(delay) * time.Minute).Format("3:04 PM"),
		Delay:            delay,
		IsDelayed:        isDelayed,
		Type:             flightType,
		CrewCount:        crewCount,
		Note:             "",
		SortTime:         arrivalTime.Add(time.Duration(delay) * time.Minute),
	}
}

func removeFlightHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	id, _ := strconv.Atoi(r.FormValue("id"))

	for i, flight := range flights {
		if flight.ID == id {
			flights = append(flights[:i], flights[i+1:]...)
			break
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func updateNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	id, _ := strconv.Atoi(r.FormValue("id"))
	note := r.FormValue("note")

	for i := range flights {
		if flights[i].ID == id {
			flights[i].Note = note
			break
		}
	}

	w.WriteHeader(http.StatusOK)
}
