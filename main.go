package main

import (
	"encoding/json"
	_"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type errorData struct {
	Num  int
	Text string
}

// defining artist information struct for unmarshling
type artistInfo struct {
	ID             int                 `json:"id"`
	Image          string              `json:"image"`
	Name           string              `json:"name"`
	Members        []string            `json:"members"`
	CreationDate   int                 `json:"creationDate"`
	FirstAlbum     string              `json:"firstAlbum"`
	Concerts       string              `json:"concertDates"`
	Locations      string              `json:"locations"`
	Relations      string              `json:"Relations"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

// struct to deal with info
// cannot be unmarshelled in one step hence why there are two structs
type artistInfo2 struct {
	ID             int                 `json:"id"`
	Image          string              `json:"image"`
	Name           string              `json:"name"`
	Members        []string            `json:"members"`
	CreationDate   int                 `json:"creationDate"`
	FirstAlbum     string              `json:"firstAlbum"`
	Concerts       []string            `json:"concertDates"`
	Locations      []string            `json:"locations"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type date struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type relation struct {
	Id             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

// removes the word index from date list
type dateremoveindex struct {
	Index []date `json:"index"`
}

type locationremoveindex struct {
	Index []location `json:"index"`
}

type relationremoveindex struct {
	Index []relation `json:"index"`
}

// let's declare a global Articles array
// that we can then populate in our main function
// to simulate a database
var (
	FullArtistInfo    []artistInfo2
	Articles          []artistInfo
	locations         []location
	Dates             []date
	DatesLocations    []relation
	DateswithIndex    dateremoveindex
	RelationwithIndex relationremoveindex
	LocationWithIndex locationremoveindex
)

var (
	eventLocations locationremoveindex
	eventDates     dateremoveindex
	Artists        []artistInfo
	dataLocation   relationremoveindex
)

// Entry point to our REST server
func main() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/response", response)
	// http.HandleFunc("/search", searchBar)

	// declaring variables for fucntions to be called
	Articles = gatherDataUp()
	Dates = gatherDataUpDates()
	locations = gatherDataUpLocations()
	DatesLocations = gatherDataUpRelation()
	appendalldata()

	// Serves and targets files that are within folder
	fs := http.FileServer((http.Dir("ui")))
	http.Handle("/ui/html/", http.StripPrefix("/ui/", fs))
//serves and target css files
	fss := http.FileServer((http.Dir("ui")))
	http.Handle("/ui/styles/", http.StripPrefix("/ui/", fss))

	//serving images
	img := http.FileServer((http.Dir("ui")))
	http.Handle("/ui/images/", http.StripPrefix("/ui/", img))
	// starting our server
	log.Println("Starting Server :http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}

// handles all request to root URL
func homePage(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("ui/html/index.html")
	if err != nil {
		errorPage(w, r, http.StatusInternalServerError)
		return
	}
	// if not homepage then return error status
	if r.URL.Path != "/" {
		errorPage(w, r, http.StatusNotFound)
		return
	}
	// ranging through struct artistInfo to gather data withins
	var result []artistInfo2
	result = FullArtistInfo
	if err := temp.Execute(w, result); err != nil {
		log.Printf("Execute Error: %v", err)
		http.Error(w, "Error when Executing", http.StatusInternalServerError)
		return
	}
	// fmt.Println(result)
}

func response(w http.ResponseWriter, r *http.Request) {
    temp, er := template.ParseFiles("ui/html/response.html")
    if er != nil {
        log.Fatal(er)
        errorPage(w, r, http.StatusInternalServerError)
        return
    }

    idValue := r.FormValue("id")

    idValueInt, err := strconv.Atoi(idValue)
    if err != nil {
       // log.Println("Invalid artist ID")
        errorPage(w, r, http.StatusBadRequest)
        return
    }

    artist := getArtistID(idValueInt)

    if artist.ID == 0 {
        // Artist not found, render the error page
        temp, er = template.ParseFiles("ui/html/error.html")
        if er != nil {
            log.Fatal(er)
            errorPage(w, r, http.StatusInternalServerError)
            return
        }
        errData := errorData{Num: http.StatusNotFound, Text: "No artist info found"}
        w.WriteHeader(http.StatusNotFound)
        temp.Execute(w, errData)
        return
    }

    temp.Execute(w, artist)
}
// functions to check & deal with all errors
func errorPage(w http.ResponseWriter, r *http.Request, errorPage int) {
	temp, er := template.ParseFiles("ui/html/error.html")
	if er != nil {
		log.Fatal(er)
		return
	}
	w.WriteHeader(errorPage)
	errData := errorData{Num: errorPage}
	if errorPage == 404 {
		errData.Text = "Page Not Found"
	} else if errorPage == 400 {
		errData.Text = "Bad Request"
	} else if errorPage == 500 {
		errData.Text = "Internal Server Error"
	}
	temp.Execute(w, errData)
}

func gatherDataUp() []artistInfo {
	data1 := myArtist() // data1 has a value of getData function
	// slices the artistData from any value, assign that to artist
	er := json.Unmarshal(data1, &Artists) // convert the data1 and pointed artists through json into a variable
	if er != nil {
		log.Fatal(er)
		return nil
	}
	for i := 0; i < len(Artists); i++ { // iterates through the length of artist
		r := artistInfo{}         // assigns any typed relation to r
		json.Unmarshal(data1, &r) // converts sliced artist and relation? and pointed r variable thr json to
		// links artists to concerts that's pointed from r
	}
	for i := 0; i < 52; i++ {
		// fmt.Println(Artists[i].Name)
	}
	return Artists // return artist result
}

func gatherDataUpRelation() []relation {
	data1 := myDatesLocations()

	er := json.Unmarshal(data1, &dataLocation)
	if er != nil {
		log.Fatal(er)
		return nil
	}
	for i := 0; i < len(dataLocation.Index); i++ {
		r := relation{}
		json.Unmarshal(data1, &r)
	}
	for i := 0; i < 52; i++ {
		// fmt.Println(dataLocation.Index[i])
	}
	return dataLocation.Index
}

func gatherDataUpDates() []date {
	data1 := myDates()

	er := json.Unmarshal(data1, &eventDates)
	if er != nil {
		log.Fatal(er)
		return nil
	}
	for i := 0; i < len(eventDates.Index); i++ {
		// r := .dateindex
		json.Unmarshal(data1, &eventDates)
	}

	// fmt.Println(eventDates.Index)

	return eventDates.Index
}

func gatherDataUpLocations() []location {
	data1 := myLocations()

	er := json.Unmarshal(data1, &eventLocations)
	if er != nil {
		log.Fatal(er)
		return nil
	}
	for i := 0; i < len(eventLocations.Index); i++ {
		// r := .dateindex
		json.Unmarshal(data1, &eventLocations)
	}

	// fmt.Println(eventLocations.Index)

	return eventLocations.Index
}

func myArtist() []byte {
	data1, e1 := http.Get("https://groupietrackers.herokuapp.com/api/artists") // requests data from the link server
	if e1 != nil {
		log.Fatal(e1)
		return nil // error reading
	}
	defer data1.Body.Close()
	data2, e2 := ioutil.ReadAll(data1.Body) // read the data from data1

	if e2 != nil {
		log.Fatal(e2)
		return nil // return the error
	}
	return data2 // else return the read data that had been requested
}

func myDates() []byte {
	data1, e1 := http.Get("https://groupietrackers.herokuapp.com/api/dates")
	if e1 != nil {
		log.Fatal(e1)
		return nil
	}
	defer data1.Body.Close()
	data2, e2 := ioutil.ReadAll(data1.Body)

	if e2 != nil {
		log.Fatal(e2)
		return nil
	}
	return data2
}

func myLocations() []byte {
	data1, e1 := http.Get("https://groupietrackers.herokuapp.com/api/locations")
	if e1 != nil {
		log.Fatal(e1)
		return nil
	}
	defer data1.Body.Close()
	data2, e2 := ioutil.ReadAll(data1.Body)

	if e2 != nil {
		log.Fatal(e2)
		return nil
	}
	return data2
}

func myDatesLocations() []byte {
	data1, e1 := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	if e1 != nil {
		log.Fatal(e1)
		return nil
	}
	defer data1.Body.Close()
	data2, e2 := ioutil.ReadAll(data1.Body)

	if e2 != nil {
		log.Fatal(e2)
		return nil
	}
	return data2
}

// appends all info in our 2nd struct into one variable
func appendalldata() []artistInfo2 {
	for i := range Articles {
		var appendingdata artistInfo2
		appendingdata.ID = i + 1
		appendingdata.Name = Artists[i].Name
		appendingdata.Members = Artists[i].Members
		appendingdata.Image = Artists[i].Image
		appendingdata.CreationDate = Artists[i].CreationDate
		appendingdata.FirstAlbum = Artists[i].FirstAlbum
		appendingdata.DatesLocations = dataLocation.Index[i].DatesLocations
		appendingdata.Locations = eventLocations.Index[i].Locations
		appendingdata.Concerts = eventDates.Index[i].Dates
		FullArtistInfo = append(FullArtistInfo, appendingdata)

	}
	// fmt.Println("----------------", FullArtistInfo[0])
	return FullArtistInfo
}

func getArtistID(id int) artistInfo2 {
	var artistnotfound artistInfo2

	for _, item := range FullArtistInfo {
		if item.ID == id {
			return item
		}
	}
	return artistnotfound
}







