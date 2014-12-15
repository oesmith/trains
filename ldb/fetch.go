package ldb

import (
	"bytes"
	"encoding/xml"
	"flag"
	"io/ioutil"
	"net/http"
)

// There's only ever going to be one access_token; use a flag.
var accessToken = flag.String("ldb_token", "",
		"Access token for the National Rail live departures board service")
var serviceAddress = flag.String("ldb_address", "",
		"Address of the National Rail live departures board service")

func NewGetDepartureBoardRequest(rows int, from string, to string) (*http.Request, error) {
	envelope := &Envelope{
		EncodingStyle: "http://schemas.xmlsoap.org/soap/encoding",
		Header: &EnvelopeHeader{
			AccessToken: &Token{
				TokenValue: *accessToken,
			},
		},
		Body: &EnvelopeBody{
			GetDepartureBoardRequest: &GetDepartureBoardRequest{
				NumRows: rows,
				CRS: from,
				FilterCRS: to,
			},
		},
	}
	return newRequest(envelope, "http://thalesgroup.com/RTTI/2012-01-13/ldb/GetDepartureBoard")
}

func ParseGetDepartureBoardResponse(resp *http.Response) (*GetDepartureBoardResponse, error) {
	if envelope, err := parseEnvelope(resp); err != nil {
		return nil, err
	} else {
		return envelope.Body.GetDepartureBoardResponse, nil
	}
}

func newRequest(envelope *Envelope, action string) (*http.Request, error) {
	bytes := &bytes.Buffer{}
	encoder := xml.NewEncoder(bytes)
	if err := encoder.Encode(envelope); err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", *serviceAddress, bytes)
	if err != nil {
		return nil, err
	}
	req.Header.Add("SOAPAction", action)
	req.Header.Add("Content-Type", "text/xml")
	return req, nil
}

func parseEnvelope(resp *http.Response) (*Envelope, error) {
	envelope := &Envelope{}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = xml.Unmarshal(body, envelope)
	if err != nil {
		return nil, err
	}
	return envelope, nil
}
