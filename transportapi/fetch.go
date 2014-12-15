package transportapi

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	urlTemplate = "http://transportapi.com/v3/uk/train/station/%s/live.json?app_id=%s&api_key=%s"
)

var (
	appId = flag.String("transportapi_app_id", "", "Transport API app ID")
	apiKey = flag.String("transportapi_api_key", "", "Transport API key")
)

func NewDeparturesRequest(from string) (*http.Request, error) {
	url := fmt.Sprintf(urlTemplate, from, *appId, *apiKey)
	return http.NewRequest("GET", url, nil)
}

func ParseResponse(resp *http.Response) (*Response, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	data := &Response{}
	err = json.Unmarshal(body, data)
	if err != nil {
		return nil, err
	}
	if data.Error != "" {
		return nil, errors.New(data.Error)
	}
	return data, nil
}
