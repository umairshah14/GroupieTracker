package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

const api string = "https://groupietrackers.herokuapp.com/api/"

func GetBody(s string) []byte {
	link := api + s
	resp, err := http.Get(link)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return body
}

func process(w http.ResponseWriter, r *http.Request) {

	var locationsunm Locs
	locBody := GetBody("locations")
	json.Unmarshal(locBody, &locationsunm)

	artBody := GetBody("artists")
	var artistunm []Artist
	json.Unmarshal(artBody, &artistunm)

	datBody := GetBody("dates")
	datesunm := Dats{}
	json.Unmarshal(datBody, &datesunm)

	relBody := GetBody("relation")
	relationsunm := Rel{}
	json.Unmarshal(relBody, &relationsunm)

	switch r.Method {
	case "GET":
		all := map[string]interface{}{
			"artists":   artistunm,
			"locations": locationsunm,
			"relations": relationsunm,
			"dates":     datesunm,
		}
		t := template.Must(template.ParseFiles(("groupietracker.html")))

		t.Execute(w, all)

	case "POST":
		ID := r.FormValue("ID")
		numID, _ := strconv.Atoi(ID)

		locs := locationsunm.Index[numID-1]
		rels := relationsunm.Index[numID-1].DatsLocs

		var dates = make(map[int][]string)

		for i, c := range locs.Locat {
			s := rels[c]
			dates[i] = s
		}

		all := Artistinfo{Artist: artistunm[numID-1], Location: locationsunm.Index[numID-1], Dates: dates, Relation: relationsunm.Index[numID-1]}

		t := template.Must(template.ParseFiles("artistinfo.html"))

		t.Execute(w, all)

	default:
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)

	}

}

func main() {
	fs := http.FileServer(http.Dir("./static"))

	http.Handle("/static/", http.StripPrefix("/static/", fs)) // handling the CSS
	// r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	http.HandleFunc("/", process)
	http.HandleFunc("/artistinfo", process)

	fmt.Printf("Starting server at port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type Locs struct {
	Index []Locations `json:"index"`
}
type Locations struct {
	ID    int      `json:"id"`
	Locat []string `json:"locations"`
	Data  string   `json:"dates"`
}

type Dats struct {
	Index []Dates `json:"index"`
}
type Dates struct {
	ID  int      `json:"id"`
	Dat []string `json:"dates"`
}

type Rel struct {
	Index []Relation `json:"index"`
}
type Relation struct {
	ID       int                 `json:"id"`
	DatsLocs map[string][]string `json:"datesLocations"`
}
type Artistinfo struct {
	Artist   interface{}
	Location interface{}
	Dates    interface{}
	Relation interface{}
}
