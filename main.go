package main

import (
  "encoding/json"
  "fmt"
  "io"
  "net/http"
)

type Weather struct {
    Location struct {
        Name    string `json:"name"`
        Country string `json:"country"`
    } `json:"location"`

    Current struct {
    TempC float64 `json:"temp_c"`
        Condition struct {
            Text string `json:"text"`
        } `json:"condition"`
    } `json:"current"`

Forecast struct {
        ForecastDay []struct {
            Hour []struct {
                TimeEpoch    int64   `json:"time_epoch"`
                TempC        float64 `json:"temp_c"`
                Condition    struct {
                    Text string `json:"text"`
                } `json:"condition"`
                ChanceOfRain float64 `json:"chance_of_rain"`
            } `json:"hour"`
        } `json:"forecastday"`
    } `json:"forecast"`
}

func main() {
  
  res, err := http.Get("https://api.openweathermap.org/data/2.5/weather?q=London,uk&appid=8a9e98d41de585beb8405200c2b50dee")
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

  var weather Weather
  err = json.Unmarshal(body, &weather)

  if err != nil {
    panic(err)
  }
  fmt.Println(weather)
  
}
