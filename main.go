package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"io"
	"fmt"
)

type Artist struct {
	ID          int      `json:"id"`
	Image       string   `json:"image"`
	Name        string   `json:"name"`
	Members     []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum  string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type ArtistsData []Artist

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch data about artists from the GroupieTrackers API
	apiURL := "https://groupietrackers.herokuapp.com/api/artists"
	resp, err := http.Get(apiURL)
	if err != nil {
		http.Error(w, "Error fetching data from API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Parse the JSON data into an ArtistsData slice
	var artistsData ArtistsData
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading data from the API", http.StatusInternalServerError)
		return
	}
	if err := json.Unmarshal(data, &artistsData); err != nil {
		http.Error(w, "Error parsing JSON data", http.StatusInternalServerError)
		return
	}

	// Parse the HTML template from the 'artists.html' file
	tmpl, err := template.ParseFiles("artists.html")
	if err != nil {
		http.Error(w, "Error parsing HTML template", http.StatusInternalServerError)
		return
	}

	// Execute the template and write it to the response
	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, artistsData); err != nil {
		http.Error(w, "Error executing HTML template", http.StatusInternalServerError)
	}
}

func main() {
	// Define a route to handle the home page
	http.HandleFunc("/", homeHandler)

	// Start the HTTP server
	fmt.Println("HTTP SERVER RUNNING AT: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
