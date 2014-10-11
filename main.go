package main

import (
	"encoding/json"
	"fmt"
	"github.com/grindhold/gominatim"
	"github.com/op/go-logging"
	"os"
	"strings"
)

var log = logging.MustGetLogger("places")
var format = "%{level} ▶ %{message}"

func main() {
	logging.SetFormatter(logging.MustStringFormatter(format))

	gominatim.SetServer("https://nominatim.openstreetmap.org/")

	if len(os.Args) < 2 {
		fmt.Println("Usage: octosquare <location>")
		os.Exit(1)

	}
	args := strings.Join(os.Args[1:], " ")

	//Get by a Querystring
	log.Info("Searching for '" + args + "'...")
	qry := gominatim.SearchQuery{
		Q: args,
	}
	resp, _ := qry.Get() // Returns []gominatim.Result
	if len(resp) > 0 {
		log.Info("Found locations:\n")
		for _, loc := range resp {
			log.Info("  -> %s (%s, %s)\n",
				loc.DisplayName, loc.Lat, loc.Lon)
		}
	} else {
		log.Error("No locations found for: '%s'", args)

	}

	places_json := os.Getenv("HOME") + "/Work/github/places/places.v1.json"
	places := New(places_json)
	places.AddPlace(args)

	file, err := os.Create(places_json)
	if err != nil {
		panic("Error opening places file for writing")
	}
	defer file.Close()

	str, err := json.MarshalIndent(places.GeoJ, "", "  ")
	if err != nil {
		panic("Error opening places file for writing")
	}
	file.Write(str)
}