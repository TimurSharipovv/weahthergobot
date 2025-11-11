package openweather

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type OpenWeatherClient struct {
	apiKey string
}

func New(apiKey string) *OpenWeatherClient {
	return &OpenWeatherClient{
		apiKey: apiKey,
	}
}

func (o OpenWeatherClient) Coordinates(city string) (Coordinates, error) {
	url := "http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=5&appid=%s"
	resp, err := http.Get((fmt.Sprintf(url, city, o.apiKey)))
	if err != nil {
		return Coordinates{}, fmt.Errorf("error ger coordinates: %w", err)
	}

	if resp.StatusCode != 200 {
		return Coordinates{}, fmt.Errorf("error getting coordinates: %d", resp.StatusCode)
	}

	var coordinatesResponse []CoordinatesResponse
	err = json.NewDecoder(resp.Body).Decode(&coordinatesResponse)
	if err != nil {
		return Coordinates{}, fmt.Errorf("error unmarshal response: %w", err)
	}

	if len(coordinatesResponse) == 0 {
		return Coordinates{}, fmt.Errorf("coordinates are empty")
	}

	return Coordinates{
		Lat: coordinatesResponse[0].Lat,
		Lon: coordinatesResponse[0].Lon,
	}, nil
}

func (o OpenWeatherClient) Weather(lat, lon float64) (Weather, error) {
	url := "https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s"
	resp, err := http.Get(fmt.Sprintf(url, lat, lon, o.apiKey))

	if err != nil {
		return Weather{}, fmt.Errorf("error getting weather: %w", err)
	}

	if resp.StatusCode != 200 {
		return Weather{}, fmt.Errorf("error fail get weather: %d", resp.StatusCode)
	}

	var weatherResponse WeatherResponse
	err = json.NewDecoder(resp.Body).Decode(&weatherResponse)

	if err != nil {
		return Weather{}, fmt.Errorf("error unmarshal weather response: %w", err)
	}

	return Weather{
		Temp: weatherResponse.Main.Temp,
	}, nil
}
