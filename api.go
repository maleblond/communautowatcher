package communautowatcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

const carAvailabilityURL = "https://www.reservauto.net/Scripts/Client/Ajax/PublicCall/Get_Car_DisponibilityJSON.asp"
const flexCarAvailabilityURL = "https://www.reservauto.net/WCF/LSI/LSIBookingServiceV3.svc/GetAvailableVehicles"

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

type flexCarAvailabilitiesResp struct {
	Data flexVehiclesResp `json:"d"`
}

type flexVehiclesResp struct {
	Vehicles []FlexCarAvailabilityResp `json:"Vehicles"`
}

type FlexCarAvailabilityResp struct {
	IsPromo   bool    `json:"isPromo"`
	CarBrand  string  `json:"CarBrand"`
	CarModel  string  `json:"CarModel"`
	CarNo     int     `json:"CarNo"`
	CarPlate  string  `json:"CarPlate"`
	Latitude  float64 `json:"Latitude"`
	Longitude float64 `json:"Longitude"`
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

func GetAvailableFlexCars(ctx context.Context, query CarQuery) ([]Car, error) {
	req, err := http.NewRequest(http.MethodGet, flexCarAvailabilityURL, nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	req.WithContext(ctx)
	q := req.URL.Query()
	q.Add("BranchID", query.BranchID)
	q.Add("LanguageID", query.LanguageID)
	q.Add("CityID", query.CityID)
	req.URL.RawQuery = q.Encode()

	carsResponse := flexCarAvailabilitiesResp{}
	res, err := http.DefaultClient.Do(req)
	if err = json.NewDecoder(res.Body).Decode(&carsResponse); err != nil {
		return nil, err
	}

	cars := []Car{}
	for _, vehicle := range carsResponse.Data.Vehicles {

		cars = append(cars, Car{
			Latitude:  vehicle.Latitude,
			Longitude: vehicle.Longitude,
			IsPromo:   vehicle.IsPromo,
			CarBrand:  vehicle.CarBrand,
			CarModel:  vehicle.CarModel,
			CarNo:     vehicle.CarNo,
			CarPlate:  vehicle.CarPlate,
		})
	}

	return cars, nil
}

func GetAvailableCars(ctx context.Context, query CarQuery) ([]Car, error) {
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
