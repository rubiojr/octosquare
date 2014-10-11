package main

import (
	"encoding/json"
	"fmt"
	"github.com/grindhold/gominatim"
	"github.com/op/go-logging"
	"github.com/rubiojr/freegeoip-client"
	"github.com/rubiojr/kingpin"
	"os"
	"strings"
)

var log = logging.MustGetLogger("places")
var format = "%{level} ▶ %{message}"

var (
	debug = kingpin.Flag("debug", "Enable debug mode.").Bool()

	setup     = kingpin.Command("setup", "Setup wizard.")
	searchCmd = kingpin.Command("search", "Search for a place.")
	//postImage   = post.Flag("image", "Image to post.").File()
	searchPlace = searchCmd.Arg("place", "Place name.").Required().Strings()

	addCmd        = kingpin.Command("add", "Add a place.")
	addCmdPlace   = addCmd.Arg("place", "Place name.").Required().String()
	addCurrentCmd = kingpin.Command("add-current",
		"Add your current position.")

	locateMeCmd = kingpin.Command("locate-me", "Print my location.")
)

func search(args string) []map[string]string {
	gominatim.SetServer("https://nominatim.openstreetmap.org/")

	if len(os.Args) < 2 {
		fmt.Println("Usage: octosquare <location>")
		os.Exit(1)

	}

	//Get by a Querystring
	log.Info("Searching for '" + args + "'...")
	qry := gominatim.SearchQuery{
		Q: args,
	}
	resp, _ := qry.Get() // Returns []gominatim.Result
	var locs []map[string]string
	if len(resp) > 0 {
		log.Info("Found locations:\n")
		for _, loc := range resp {
			p := map[string]string{
				"name": loc.DisplayName,
				"lat":  loc.Lat,
				"lon":  loc.Lon,
			}
			locs = append(locs, p)
			log.Info("  -> %s (%s, %s)\n",
				loc.DisplayName, loc.Lat, loc.Lon)
		}
	} else {
		log.Error("No locations found for: '%s'", args)

	}

	return locs
}

func add(args string) {
	gominatim.SetServer("https://nominatim.openstreetmap.org/")

	if len(os.Args) < 2 {
		fmt.Println("Usage: octosquare <location>")
		os.Exit(1)

	}

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

func main() {

	switch kingpin.Parse() {
	// Register user
	case "setup":
		println("TODO")
	case "add":
		add(*addCmdPlace)
	case "add-current":
		loc, err := freegeoip_client.GetLocation()
		if err != nil {
			fmt.Println("Error getting my location.")
		}
		add(loc.CountryName + ", " + loc.City)
	case "search":
		search(strings.Join(*searchPlace, " "))
	case "locate-me":
		loc, err := freegeoip_client.GetLocation()
		if err != nil {
			fmt.Println("Error getting my location.")
		}
		fmt.Printf("%s, %s (%f, %f)", loc.CountryName,
			loc.City, loc.Longitude, loc.Latitude)
	}
}