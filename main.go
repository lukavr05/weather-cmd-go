package main

import (
  // "fmt"
  // "io"
  "net/http"
)



func main() {
  
  res, err := http.Get("api.openweathermap.org/data/2.5/weather?q=London,uk&APPID=8a9e98d41de585beb8405200c2b50dee")

  if err != nil {
    panic(err)
  }

  defer res.Body.Close()

  if res.StatusCode != 200 {
    panic("Weather API not available!")
  }
  
}
