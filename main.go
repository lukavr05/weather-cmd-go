package main

import (
  "encoding/json"
  "fmt"
  "io"
  "net/http"
  "flag"
  "time"
)

type WeatherForecast struct {
  Cod     string  `json:"cod"`
	Message int     `json:"message"`
	Cnt     int     `json:"cnt"`
	List    []Entry `json:"list"`
	City    System  `json:"city"`
}

type Entry struct {
	Dt       int64        `json:"dt"`
	Main     MainWeather  `json:"main"`
	Weather  []Weather    `json:"weather"`
	Clouds   Clouds       `json:"clouds"`
	Wind     Wind         `json:"wind"`
	Visibility int        `json:"visibility"`
	Pop      float64      `json:"pop"`
	Rain     *Rain        `json:"rain,omitempty"` // optional, use pointer
	Sys      System       `json:"sys"`
	DtTxt    string       `json:"dt_txt"`
}

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

type Rain struct {
	ThreeH float64 `json:"3h"`
}

// Clouds struct for the "clouds" field
type Clouds struct {
	All int `json:"all"`
}

// System struct for the "sys" field
type System struct {
  ID        int     `json:"id"`
	Name      string  `json:"name"`
  Country   string  `json:"country"`
	Population int    `json:"population"`
	Timezone  int     `json:"timezone"`
	Sunrise   int64   `json:"sunrise"`
	Sunset    int64   `json:"sunset"`}

func formatDate(date string) string {
  parsedTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		panic(err)
	}
  formattedDate := parsedTime.Format("02/01/2006")
  return formattedDate
}

// formatting the degrees so that wea can print the cardinal direction
func formatWind(degree int) string {

  if (degree >= 337 && degree < 22) {
    return "North"
  } else if (degree >= 22 && degree < 67) {
    return "North-East"
  } else if (degree >= 67 && degree < 112) {
    return "East"
  } else if (degree >= 112 && degree < 157) {
    return "South-East"
  } else if (degree >= 157 && degree < 202) {
    return "South"
  } else if (degree >= 202 && degree < 247) {
    return "South-West"
  } else if (degree >= 247 && degree < 292) {
    return "West"
  } else {
    return "North West"
  } 
}


// printing the default weather report
func printMain(w WeatherResponse) {
  fmt.Printf("City:            %s\n", w.Name)
  fmt.Printf("Description:     %s\n", w.Weather[0].Description)
  fmt.Printf("Temperature:     %.1f°C\n", w.Main.Temp)
  fmt.Printf("Feels like:      %.1f°C\n", w.Main.FeelsLike)
}

// printing a detailed weather report
func printExtended(w WeatherResponse) {
  printMain(w)
  fmt.Printf("Max Temperature: %.1f°C\n", w.Main.TempMax)
  fmt.Printf("Min Temperature: %.1f°C\n", w.Main.TempMin)
  fmt.Printf("Wind Speed:      %.1fm/s\n", w.Wind.Speed)
  fmt.Printf("Wind Direction:  %d° (%s)\n", w.Wind.Deg, formatWind(w.Wind.Deg))
  fmt.Printf("Humidity:        %d%%\n", w.Main.Humidity)
  fmt.Printf("Visibility:      %dm\n", w.Visibility)
  fmt.Printf("Cloud Coverage:  %d%%\n",w.Clouds.All)
}

func printMainForecast(wf WeatherForecast) {
  today := wf.List[0].DtTxt[:10]
  today_formatted := formatDate(today)
  fmt.Printf("\nWeather Forecast for today (%s) in %s\n",today_formatted, wf.City.Name)
  for _, entry := range wf.List {
    currentDate := entry.DtTxt[:10]
    time := entry.DtTxt[11:16]
  
    if currentDate != today {
      break
    }
    
    fmt.Printf("%s ===========================\n", time)
    fmt.Printf("  Weather:         %s\n", entry.Weather[0].Description)
    fmt.Printf("  Temperature:     %.1f°C\n",entry.Main.Temp)
  }
}

func printExtendedForecast(wf WeatherForecast) {
  fmt.Printf("Extended Weather Forecast for %s:\n", wf.City.Name)
    
  var lastDate string
  for _, entry := range wf.List {
    date := entry.DtTxt[:10] // Extract the date (YYYY-MM-DD)
    time := entry.DtTxt[11:16] // Extract the time (HH:MM)

    // Print the date only if it has changed
    if date != lastDate {
      // Format the date
      formattedDate := formatDate(date)
      fmt.Printf("\nDate: %s ===================================\n", formattedDate)
      lastDate = date // Update the last date to the current date
    }

    // Print weather information for the current entry
    fmt.Printf("\n  %s\n", time)
    fmt.Printf("    Weather: %s\n", entry.Weather[0].Description)
    fmt.Printf("    Temperature: %.1f°C\n", entry.Main.Temp)
    fmt.Printf("    Humidity: %d%%\n", entry.Main.Humidity)
    fmt.Printf("    Wind Speed: %.1fm/s, Direction: %d° (%s)\n", entry.Wind.Speed, entry.Wind.Deg, formatWind(entry.Wind.Deg))
  }
}

func main() {

  // defining our flags
  extendPtr := flag.Bool("e", false, "show an extended view of the weather report or forecast")
  forecastPtr := flag.Bool("f", false, "show the weather forecast for the next 5 days")

  flag.Parse()
  var city string

  // check if we have any arguments
  if len(flag.Args()) >= 1 {
    city = flag.Args()[0]
  } else {
    // default case
    city = "London"
  }
  
  var APIcall1 string

  // updating the API call depending on if the -f flag has been set
  if (!*forecastPtr) {
    APIcall1 = "https://api.openweathermap.org/data/2.5/weather?q="
  } else {
    APIcall1 = "https://api.openweathermap.org/data/2.5/forecast?q="
  }

  // defining the rest of the API call and making it
  APIcall2 := ",uk&appid=8a9e98d41de585beb8405200c2b50dee&units=metric"
  res, err := http.Get(APIcall1 + city + APIcall2)
  
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

  var weatherf WeatherForecast
  var weather WeatherResponse

  if (*forecastPtr) {
    err = json.Unmarshal(body, &weatherf)
  } else {
    err = json.Unmarshal(body, &weather)
  }

  if err != nil {
    panic(err)
  }

  if (*extendPtr && !*forecastPtr) {
    printExtended(weather)
  } else if (!*extendPtr && !*forecastPtr) {
    printMain(weather)
  } else if (*forecastPtr && !*extendPtr) {
    printMainForecast(weatherf)
  } else {
    printExtendedForecast(weatherf)
  }
    
}

