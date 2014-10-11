package main

import (
	"encoding/json"
	"github.com/grindhold/gominatim"
	"os"
	"strconv"
)

type GeoJ struct {
	Features []Feature `json:"features"`
	Type     string    `json:"type"`
}

type Feature struct {
	Geometry   Geometry   `json:"geometry"`
	Type       string     `json:"type"`
	Properties Properties `json:"properties"`
}

type Properties struct {
	Title string
}

type Geometry struct {
	Coordinates [2]float64 `json:"coordinates"`
	Type        string     `json:"type"`
}

type Places struct {
	GeoJ GeoJ
}

func New(file string) *Places {
	p := new(Places)
	p.GeoJ = GeoJ{
		Features: []Feature{},
		Type:     "FeatureCollection",
	}

	gominatim.SetServer("https://nominatim.openstreetmap.org/")

	handle, err := os.Open(file) // For read access.
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	enc := json.NewDecoder(handle)
	enc.Decode(&p.GeoJ)

	return p
}

func (places *Places) AddPlace(place string) {

	qry := gominatim.SearchQuery{
		Q: place,
	}
	resp, _ := qry.Get() // Returns []gominatim.Result

	if len(resp) > 0 {
		loc := resp[0]
		lat, err := strconv.ParseFloat(loc.Lat, 64)
		if err != nil {
			panic("Error parsing float.")
		}
		lon, err := strconv.ParseFloat(loc.Lon, 64)
		if err != nil {
			panic("Error parsing float.")
		}
		geom := Geometry{
			Coordinates: [2]float64{lon, lat},
			Type:        "Point",
		}
		feat := Feature{
			Geometry: geom,
			Type:     "Feature",
			Properties: Properties{
				Title: loc.DisplayName,
			},
		}
		places.GeoJ.Features = append(places.GeoJ.Features, feat)
	} else {
		log.Error("No locations found for: '%s'", place)
	}
}