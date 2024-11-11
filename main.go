package main 

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// storing types for the weather forecast
type WeatherForecast struct {
	Cod     string  `json:"cod"`
	Message int     `json:"message"`
	Cnt     int     `json:"cnt"`
	List    []Entry `json:"list"`
	City    System  `json:"city"`
}

// storing the values of each entry of a forecast
type Entry struct {
	Dt         int64       `json:"dt"`
	Main       MainWeather `json:"main"`
	Weather    []Weather   `json:"weather"`
	Clouds     Clouds      `json:"clouds"`
	Wind       Wind        `json:"wind"`
	Visibility int         `json:"visibility"`
	Pop        float64     `json:"pop"`
	Rain       *Rain       `json:"rain,omitempty"` // optional, use pointer
	Sys        System      `json:"sys"`
	DtTxt      string      `json:"dt_txt"`
}

// storing weather response for the current weather
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
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  int     `json:"pressure"`
	Humidity  int     `json:"humidity"`
	SeaLevel  int     `json:"sea_level"`
	GrndLevel int     `json:"grnd_level"`
}

// Wind struct for the "wind" field
type Wind struct {
	Speed float64 `json:"speed"`
	Deg   int     `json:"deg"`
}

// rain struct for the "rain" field
type Rain struct {
	ThreeH float64 `json:"3h"`
}

// Clouds struct for the "clouds" field
type Clouds struct {
	All int `json:"all"`
}

// System struct for the "sys" field
type System struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Country    string `json:"country"`
	Population int    `json:"population"`
	Timezone   int    `json:"timezone"`
	Sunrise    int64  `json:"sunrise"`
	Sunset     int64  `json:"sunset"`
}

type Config struct {
	DefaultCity string `json:default_city`
}

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

// function for formatting the date in YYYY-MM-DD into DD/MM/YYYY
func formatDate(date string) string {
	parsedTime := Must(time.Parse("2006-01-02", date))
	formattedDate := parsedTime.Format("02/01/2006")
	return formattedDate
}

// formatting the degrees so that wea can print the cardinal direction
func formatWind(degree int) string {
	if degree >= 337 && degree < 22 {
		return "North"
	} else if degree >= 22 && degree < 67 {
		return "North-East"
	} else if degree >= 67 && degree < 112 {
		return "East"
	} else if degree >= 112 && degree < 157 {
		return "South-East"
	} else if degree >= 157 && degree < 202 {
		return "South"
	} else if degree >= 202 && degree < 247 {
		return "South-West"
	} else if degree >= 247 && degree < 292 {
		return "West"
	} else {
		return "North West"
	}
}

// printing the default weather report
func printMain(w WeatherResponse) {
	fmt.Printf("\nCity:              %s\n", w.Name)
	fmt.Printf("Description:       %s\n", w.Weather[0].Description)
	fmt.Printf("Temperature:       %.1f°C\n", w.Main.Temp)
	fmt.Printf("Feels like:        %.1f°C\n", w.Main.FeelsLike)
}

// printing a detailed weather report
func printExtended(w WeatherResponse) {
	fmt.Printf("\nCity:              %s\n", w.Name)
	fmt.Printf("Description:       %s\n", w.Weather[0].Description)
	fmt.Printf(
		"Temperature:       %.1f°C\t\tWind Speed:        %.1fm/s\n",
		w.Main.Temp,
		w.Wind.Speed,
	)
	fmt.Printf(
		"Feels like:        %.1f°C\t\tWind Direction:    %d° (%s)\n",
		w.Main.FeelsLike,
		w.Wind.Deg,
		formatWind(w.Wind.Deg),
	)
	fmt.Printf(
		"Max Temperature:   %.1f°C\t\tVisibility:        %dm\n",
		w.Main.TempMax,
		w.Visibility,
	)
	fmt.Printf(
		"Min Temperature:   %.1f°C\t\tHumidity:          %d%%\n",
		w.Main.TempMin,
		w.Main.Humidity,
	)
	fmt.Printf(
		"Cloud Coverage:    %d%%\t\t\tGround Pressure:   %dhPa\n",
		w.Clouds.All,
		w.Main.GrndLevel,
	)
}

// printing the default forecast of the current day
func printMainForecast(wf WeatherForecast) {
	// extract the date from the first entry and format it
	today := wf.List[0].DtTxt[:10]
	today_formatted := formatDate(today)

	// print current weather
	fmt.Printf("\nWeather Forecast for today (%s) in %s\n", today_formatted, wf.City.Name)

	// iterate through the list and print weather for each time until the date is not today
	for _, entry := range wf.List {
		currentDate := entry.DtTxt[:10]
		time := entry.DtTxt[11:16]

		if currentDate != today {
			break
		}

		fmt.Printf("%s =================================\n", time)
		fmt.Printf("  Weather:         %s\n", entry.Weather[0].Description)
		fmt.Printf("  Temperature:     %.1f°C\n", entry.Main.Temp)
	}
}

// print an extended version of the weather forecast for the specified number of days
func printExtendedForecast(wf WeatherForecast, days int) {
	// Validating number of days input by the user
	if days > 5 {
		days = 5
	}

	if days <= 0 {
		days = 1
	}

	fmt.Printf("\nExtended Weather Forecast for %s for the next %d day(s):\n", wf.City.Name, days)

	// store the last date
	var lastDate string
	var printedDays int

	for _, entry := range wf.List {
		date := entry.DtTxt[:10]   // Extract the date (YYYY-MM-DD)
		time := entry.DtTxt[11:16] // Extract the time (HH:MM)

		// Print the date only if it has changed
		if date != lastDate {
			if printedDays >= days {
				break
			}

			// Format the date
			formattedDate := formatDate(date)
			fmt.Printf(
				"\n%s ============================================================================\n",
				formattedDate,
			)
			lastDate = date // Update the last date to the current date
			printedDays++
		}

		// Print weather information for the current entry
		fmt.Printf("\t%s\t\t", time)
		fmt.Printf("Weather: %s\n", entry.Weather[0].Description)
		fmt.Printf(
			"\t\t\tTemperature: %.1f°C\t\tHumidity: %d%%\n",
			entry.Main.Temp,
			entry.Main.Humidity,
		)
		fmt.Printf(
			"\t\t\tWind Speed: %.1fm/s\t\tDirection: %d° (%s)\n",
			entry.Wind.Speed,
			entry.Wind.Deg,
			formatWind(entry.Wind.Deg),
		)
		fmt.Println(
			"\t-------------------------------------------------------------------------------",
		)
	}
}

func loadConfig(path string) (*Config, error) {
	config := &Config{}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Create default config if file doesn't exist
		config.DefaultCity = "London" // Default city
		err := SaveConfig(path, config)
		if err != nil {
			return nil, err
		}
	} else {
		// Open the config file
		file := Must(os.Open(path))
		
		defer file.Close()

		// Decode the JSON into the config struct
		decoder := json.NewDecoder(file)
		err = decoder.Decode(config)
		if err != nil {
			return nil, err
		}
	}
	return config, nil
}

 func SaveConfig(path string, config *Config) error {
	file := Must(os.Create(path))
	
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(config)
}

func main() {
	configPath := "config.json"
	system_default := "London"
	dflt := ""

	config, err := loadConfig(configPath)
	if err != nil {
		fmt.Println(" !!! Error loading configuration: ", err)
		dflt = system_default
	} else {
		user_default := config.DefaultCity
		dflt = user_default
	}

	// defining flags
	extendPtr := flag.Bool("e", false, "show an extended view of the weather report or forecast")
	forecastPtr := flag.Bool("f", false, "show the weather forecast")
	numDaysPtr := flag.Int(
		"days",
		1,
		"specifies the number of days (1-5) to show for the forecast, including the current day",
	)
	setDefaultPtr := flag.String(
		"setdefault",
		dflt,
		"set the default location to set for the program",
	)
	defaultPtr := flag.Bool(
		"default",
		false,
		"show the default location for weather reporting, if no config has been made, the system defaults to 'London'",
	)

	flag.Parse()
	var city string

	if *setDefaultPtr != dflt {
		config.DefaultCity = *setDefaultPtr
		fmt.Printf("Setting default city to %s...\n", *setDefaultPtr)
		err := SaveConfig(configPath, config)
		if err != nil {
			fmt.Println(" !!! Error saving configuration: ", err)
		}
		fmt.Printf("Default city successfully set to %s!\n", config.DefaultCity)
		dflt = config.DefaultCity
	}

	// check if we have any arguments
	if len(flag.Args()) >= 1 {
		city = flag.Args()[0]
	} else {
		// default case
		city = dflt
	}

	var APIcall1 string

	// updating the API call depending on if the -f flag has been set
	if !*forecastPtr {
		APIcall1 = "https://api.openweathermap.org/data/2.5/weather?q="
	} else {
		APIcall1 = "https://api.openweathermap.org/data/2.5/forecast?q="
	}

	// defining the rest of the API call and making it
	APIcall2 := ",uk&appid=8a9e98d41de585beb8405200c2b50dee&units=metric"
	res := Must(http.Get(APIcall1 + city + APIcall2))

	// closing the connection
	defer res.Body.Close()
	if res.StatusCode != 200 {
		panic("Weather API not available! Please check the city name or try again!")
	}

	// reading the response and storing in the body
	body := Must(io.ReadAll(res.Body))

	var weatherf WeatherForecast
	var weather WeatherResponse

	if *forecastPtr {
		// if we have the -f flag, store the information in the WeatherForecast struct
		err = json.Unmarshal(body, &weatherf)
	} else {
		// if we don't have the -f flag, store the information in the WeatherResponse struct
		err = json.Unmarshal(body, &weather)
	}

	if err != nil {
		panic(err)
	}

	// checking for -e and -f flags and printing the response accordingly
	if *extendPtr && !*forecastPtr {
		printExtended(weather)
	} else if !*extendPtr && !*forecastPtr {
		printMain(weather)
	} else if *forecastPtr && !*extendPtr {
		printMainForecast(weatherf)
	} else {
		printExtendedForecast(weatherf, *numDaysPtr)

		if *numDaysPtr > 5 {
			fmt.Println("\n   !!! Specified number of days greater than maximum supplied by the API, used maximum of 5")
		}

		if *numDaysPtr <= 0 {
			fmt.Println("\n   !!! Specified number of days less than minimum required, used minimum of 1")
		}
	}

	if *defaultPtr {
		fmt.Printf("\nDefault city currently set to %s", dflt)
	}
}
