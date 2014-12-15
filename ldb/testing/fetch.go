package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/oesmith/trains/ldb"
)

var (
	from = flag.String("from", "PAD", "CRS code of origin station")
	to = flag.String("to", "", "CRS code of destination station")
	rows = flag.Int("rows", 5, "Number of rows of results to return")
)

func main() {
	flag.Parse()
	req, err := ldb.NewGetDepartureBoardRequest(*rows, *from, *to)
	if err != nil {
		log.Fatal(err)
	}
	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	data, err := ldb.ParseGetDepartureBoardResponse(res)
	if err != nil {
		log.Fatal(err)
	}
	r := data.GetStationBoardResult
	log.Printf("result: %+v", r)
	for _, s := range(r.TrainServices.Service) {
		log.Printf("origin: %+v", s.Origin.Location)
		log.Printf("destination: %+v", s.Destination.Location)
		log.Printf("service: %+v", s)
	}
}
