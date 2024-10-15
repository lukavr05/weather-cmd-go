# Command Line Weather
CMD Terminal app that provides relevant weather information in terminal using Go and the [OpenWeather](https://openweathermap.org/) API

### Installing
- clone this repo
- cd into the 'weather-cmd-go'
- build the app with `go build main.go`
    - a binary should be created without any errors.

## Usage 
- After the steps above, run the app in your terminal with `./weather`
- When running the program, add a city in the UK as an argument to specify a location (e.g. - `go run main.go {city_name}`)
- For a view of the **current** weather in a location, no flags are necessary
    - The weather-cmd-go app will show a simple view of the default city's forecast. (set to `London`)
    - For an extended view of the forecast with additional information, add flag `-e` **before** your location argument
- For a view of the **forecasted** weather in a location, add the `-f` flag **before** any arguments
    - The app will show a simple view of the forecasted weather for the current day
    - For an extended view of the location's forecast for the current day, use the `-e` flag
    - To specify the number of days to be shown, **including the current day**, use the `-days` flag followed by the number of days to show (e.g - `-days {num}`)


> [!WARNING]
> for now, **ONLY cities in England** are fetched as a location from the API.

---

### Defaults

- Cities: London is set to the default city, if no arguments provided **and** no configuration is found.
    - change this by adding `-setdefault {city_name}` **before** your arguments

- Number of Days: default forecast days is `1`, you can change this by this flag `-flag {number}`


### Flags
This program uses flags to specify the following parameters:
```    
    -f              specifies that a weather report is being requested, defaulting to the current day's forecast
    -e              specifies an extended weather report, supplying additional information to the current weather report 
                    as well as increasing the range of the weather forecast
    -days           an int specifying the number of days to display for the extended weather forecast (between 1 and 5)
    -default        displays the current default city stored in the config file
    -setdefault     updates the default value in the config file
```
you can also use `-h` or `-help` to see all the available options.

### Configuration

This program makes use of a configuration file to store the default city that is used when no specific city is given as an argument. 

- use `-setdefault` to set the default location in the config file.

> Without setting the default or having the config file in the path, the system will use "London" as its default. This can be changed in the program itself in the main function denoted by the
`system_default` variable#

> If no `config.json` file exists in the current directory, calling `-setdefault` will create one
