package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"

	"github.com/oesmith/trains/ldb"
	"github.com/oesmith/trains/transportapi"
)

var (
	client = &http.Client{}
	path = regexp.MustCompile(`^/dep/([A-Z]{3})/([A-Z]{3})$`)
	tmpl = template.Must(template.New("departures").Parse(`<!doctype html>
<html>
<head>
<meta name="viewport" content="width=device-width, initial-scale=1.0, user-scalable=yes">
<title>{{.OriginCode}} -> {{.DestinationCode}}</title>
<style>
	body { font-family: monospace }
	h1 { font-size: 1.0em; font-weight: normal }
	.error { color: red }
	.unknown { color: #aaa }
</style>
</head>
<body>
<h1>{{.OriginCode}} -> {{.DestinationCode}}</h1>

{{if .LdbError}}<p class="error">{{.LdbError.Error}}</p>{{end}}
{{if .TransportapiError}}<p class="error">{{.TransportapiError.Error}}</p>{{end}}

{{if .Departures}}
<table width="100%">
{{range .Departures}}
<tr>
<td>{{.ScheduledDeparture}}</td>
<td>{{.EstimatedDeparture}}</td>
{{if .ScheduledPlatform}}<td>{{.ScheduledPlatform}}</td>{{else}}<td class="unknown">?</td>{{end}}
{{if .EstimatedPlatform}}<td>{{.EstimatedPlatform}}</td>{{else}}<td class="unknown">?</td>{{end}}
<td>{{.DestinationCode}}</td>
<td>{{.DestinationName}}</td>
</tr>
{{end}}
</table>
{{end}}

</table>
</body>
</html>`))
)

type departure struct {
	ScheduledDeparture string
	EstimatedDeparture string
	DestinationCode string
	DestinationName string
	ScheduledPlatform string
	EstimatedPlatform string
}

func fetch(req *http.Request) (*http.Response, error) {
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Invalid status code: %d", res.StatusCode)
	}
	return res, nil
}

func fetchDepartures(from string, to string) (*ldb.GetDepartureBoardResponse, error) {
	req, err := ldb.NewGetDepartureBoardRequest(5, from, to)
	if err != nil {
		return nil, err
	}
	res, err := fetch(req)
	if err != nil {
		return nil, err
	}
	data, err := ldb.ParseGetDepartureBoardResponse(res)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func fetchLive(from string) (*transportapi.Response, error) {
	req, err := transportapi.NewDeparturesRequest(from)
	if err != nil {
		return nil, err
	}
	res, err := fetch(req)
	if err != nil {
		return nil, err
	}
	data, err := transportapi.ParseResponse(res)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func merge(l *ldb.GetDepartureBoardResponse, t *transportapi.Response) []*departure {
	deps := make([]*departure, len(l.GetStationBoardResult.TrainServices.Service))

	for i := range(deps) {
		s := l.GetStationBoardResult.TrainServices.Service[i]
		deps[i] = &departure{
			DestinationCode: s.Destination.Location.CRS,
			DestinationName: s.Destination.Location.LocationName,
			ScheduledDeparture: s.STD,
			EstimatedDeparture: s.ETD,
			EstimatedPlatform: s.Platform,
		}
	}

	// n^2, woohooo!
	for _, td := range(t.Departures.All) {
		for _, d := range(deps) {
			if td.DestinationName == d.DestinationName && td.AimedDepartureTime == d.ScheduledDeparture {
				d.ScheduledPlatform = td.Platform
			}
		}
	}

	return deps
}

func departures(w http.ResponseWriter, r *http.Request) {
	parts := path.FindStringSubmatch(r.URL.Path)
	if len(parts) != 3 {
		http.Error(w, fmt.Sprintf("Bad request %v", parts), 400)
		return
	}

	dep, depErr := fetchDepartures(parts[1], parts[2])
	live, liveErr := fetchLive(parts[1])

	var deps []*departure
	if depErr == nil && liveErr == nil {
		deps = merge(dep, live)
	}

	d := &struct{
		OriginCode string
		DestinationCode string
		Departures []*departure
		LdbError error
		TransportapiError error
	}{parts[1], parts[2], deps, depErr, liveErr}

	if err := tmpl.Execute(w, d); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func main() {
	flag.Parse()
	http.HandleFunc("/dep/", departures)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
