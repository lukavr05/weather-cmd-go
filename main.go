package main

import (
  "encoding/json"
  "fmt"
  "io"
  "net/http"
)

type WeatherResponse struct {
	Coord      Coordinates `json:"coord"`
	Weather    []Weather   `json:"weather"`
	Base       string      `json:"base"`
	Main       MainWeather `json:"main"`
	Visibility int         `json:"visibility"`
	Wind       Wind        `json:"wind"`
	Clouds     Clouds      `json:"clouds"`
	Dt         int         `json:"dt"`
	Sys        System      `json:"sys"`
	Timezone   int         `json:"timezone"`
	ID         int         `json:"id"`
	Name       string      `json:"name"`
	Cod        int         `json:"cod"`
}

// Coordinates struct for the "coord" field
type Coordinates struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

// Weather struct for the "weather" field
type Weather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// MainWeather struct for the "main" field
type MainWeather struct {
	Temp       float64 `json:"temp"`
	FeelsLike  float64 `json:"feels_like"`
	TempMin    float64 `json:"temp_min"`
	TempMax    float64 `json:"temp_max"`
	Pressure   int     `json:"pressure"`
	Humidity   int     `json:"humidity"`
	SeaLevel   int     `json:"sea_level"`
	GrndLevel  int     `json:"grnd_level"`
}

// Wind struct for the "wind" field
type Wind struct {
	Speed float64 `json:"speed"`
	Deg   int     `json:"deg"`
}

// Clouds struct for the "clouds" field
type Clouds struct {
	All int `json:"all"`
}

// System struct for the "sys" field
type System struct {
	Type     int    `json:"type"`
	ID       int    `json:"id"`
	Country  string `json:"country"`
	Sunrise  int    `json:"sunrise"`
	Sunset   int    `json:"sunset"`
}

func printMain(w WeatherResponse) {
  fmt.Printf("Location:        %s\n", w.Name)
  fmt.Printf("Time Zone:       %d\n",w.Timezone)
  fmt.Printf("Description:     %s\n", w.Weather[0].Description)
  fmt.Printf("Temperature:     %.1f°C\n", w.Main.Temp)
  fmt.Printf("Feels like:      %.1f°C\n", w.Main.FeelsLike)
  fmt.Printf("Humidity:        %d%%\n", w.Main.Humidity)
}

func main() {
  
  res, err := http.Get("https://api.openweathermap.org/data/2.5/weather?q=London,uk&appid=8a9e98d41de585beb8405200c2b50dee&units=metric")
  if err != nil {
    panic(err)
  }

  defer res.Body.Close()

  if res.StatusCode != 200 {
    panic("Weather API not available!")
  }

  body, err := io.ReadAll(res.Body)
  if err != nil {
    panic(err)
  }

  var weather WeatherResponse
  err = json.Unmarshal(body, &weather)

  if err != nil {
    panic(err)
  }
  printMain(weather)
}
