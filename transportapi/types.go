package transportapi

type Response struct {
	Arrivals *Trains `json:"arrivals"`
	Date string `json:"date"`
	Departures *Trains `json:"departures"`
	Error string `json:"error"`
	RequestTime string `json:"request_time"`
	StationCode string `json:"station_code"`
	StationName string `json:"station_name"`
	TimeOfDay string `json:"time_of_day"`
}

type Trains struct {
	All []*Train `json:"all"`
}

type Train struct {
	AimedArrivalTime string `json:"aimed_arrival_time"`
	AimedDepartureTime string `json:"aimed_departure_time"`
	AimedPassTime string `json:"aimed_pass_time"`
	BestArrivalEstimateMins int `json:"best_arrival_estimate_mins"`
	BestDepartureEstimateMins int `json:"best_departure_estimate_mins"`
	DestinationName string `json:"destination_name"`
	ExpectedArrivalTime string `json:"expected_arrival_time"`
	ExpectedDepartureTime string `json:"expected_departure_time"`
	Mode string `json:"mode"`
	Operator string `json:"operator"`
	OriginName string `json:"origin_name"`
	Platform string `json:"platform"`
	Service string `json:"service"`
	Source string `json:"source"`
	Status string `json:"status"`
	TrainUID string `json:"train_uid"`
}
