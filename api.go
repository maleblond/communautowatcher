package communautowatcher

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"gopkg.in/yaml.v2"
)

const carAvailabilityURL = "https://www.reservauto.net/Scripts/Client/Ajax/PublicCall/Get_Car_DisponibilityJSON.asp"

type carAvailabilitiesResp struct {
	Data []carAvailabilityResp `yaml:"data"`
}

type carAvailabilityResp struct {
	StationName string  `yaml:"strNomStation"`
	Distance    float64 `yaml:"Distance"`
	NbrRes      int     `yaml:"NbrRes"`
	Latitude    float64 `yaml:"Latitude"`
	Longitude   float64 `yaml:"Longitude"`
}

// The API does not return valid JSON, so we need to do some gymnastics (see ./samples/car_availabilities.txt):
func parseAvailableCarResponse(res *http.Response) (carAvailabilitiesResp, error) {
	parsedBody := carAvailabilitiesResp{}
	body, err := io.ReadAll(res.Body)

	if err != nil {
		return parsedBody, err
	}

	bodyStr := string(body)
	bodyStr = bodyStr[1 : len(bodyStr)-1] // Remove parenthesises that wraps the response body

	// Yaml requires a space between the colon and the value (i.e. `key: value`) but the API doesn't include those.
	// This is not perfect (i.e. it may alter station names that is included in the payload), but it's good enough
	// for my current use case.
	bodyStr = strings.ReplaceAll(bodyStr, ":", ": ")

	err = yaml.Unmarshal([]byte(bodyStr), &parsedBody)

	return parsedBody, err
}

func GetAvailableCars(query CarQuery) ([]Car, error) {
	res, err := http.PostForm(carAvailabilityURL,
		url.Values{
			"CurrentLanguageID": {"1"},
			"CityID":            {query.CityID},
			"StartDate":         {query.StartDate.Format("02/01/2006 15:04")},
			"EndDate":           {query.EndDate.Format("02/01/2006 15:04")},
			"Accessories":       {"0"},
			"FeeType":           {"80"},
			"Latitude":          {query.FromLatitude},
			"Longitude":         {query.FromLongitude},
		})

	if err != nil {
		return nil, err
	}

	parsedResp, err := parseAvailableCarResponse(res)

	if err != nil {
		return nil, err
	}

	cars := []Car{}

	for _, availability := range parsedResp.Data {
		if availability.NbrRes == 0 {
			cars = append(cars, Car{
				Distance:     availability.Distance,
				Latitude:     availability.Latitude,
				LocationName: availability.StationName,
				Longitude:    availability.Longitude,
			})
		}
	}

	return cars, nil
}
