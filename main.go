package main

import (
	"encoding/json"
	"html/template"
	"net/http"
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
	if err := json.NewDecoder(resp.Body).Decode(&artistsData); err != nil {
		http.Error(w, "Error parsing JSON data", http.StatusInternalServerError)
		return
	}

	// Create an HTML template to display the artists with images
	tmpl := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>GroupieTrackers Artists</title>
		<style>
		  .artist-grid {
		    display: grid;
		    grid-template-columns: repeat(3, 1fr);
		    gap: 20px;
		  }
		
		  .artist-card {
		    border: 1px solid #ccc;
		    padding: 20px;
		    text-align: center;
		  }
		
		  .artist-card img {
		    max-width: 100%;
		  }
		</style>
	</head>
	<body>
		<h1>GroupieTrackers Artists</h1>
		<div class="artist-grid">
		{{range .}}
			<div class="artist-card">
				<img src="{{.Image}}" alt="{{.Name}}" />
				<h2>{{.Name}}</h2>
				<p>Members: {{range .Members}}{{.}}, {{end}}</p>
				<p>Creation Date: {{.CreationDate}}</p>
				<p>First Album: {{.FirstAlbum}}</p>
				<p>Locations: <a href="{{.Locations}}">View Locations</a></p>
				<p>Concert Dates: <a href="{{.ConcertDates}}">View Concert Dates</a></p>
				<p>Relations: <a href="{{.Relations}}">View Relations</a></p>
			</div>
		{{end}}
		</div>
	</body>
	</html>
	`

	// Parse the HTML template and serve the content
	t, err := template.New("artists").Parse(tmpl)
	if err != nil {
		http.Error(w, "Error parsing HTML template", http.StatusInternalServerError)
		return
	}

	// Execute the template and write it to the response
	w.Header().Set("Content-Type", "text/html")
	if err := t.Execute(w, artistsData); err != nil {
		http.Error(w, "Error executing HTML template", http.StatusInternalServerError)
	}
}

func main() {
	// Define a route to handle the home page
	http.HandleFunc("/", homeHandler)

	// Start the HTTP server
	http.ListenAndServe(":8080", nil)
	
}
