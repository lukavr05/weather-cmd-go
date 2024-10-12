This is my command line program that provides relevant weather information in terminal using Go and the OpenWeather API

Providing a city in the UK as an argument will retreive the information for that city. Without an argument, the default is set to London

This program also uses flags to specify whether a forecast or extended weather report is desired
    
    -f specifies that a weather report is being requested, defaulting to the current day's forecast
    -e specifies an extended weather report, supplying additional information to the current weather report as well as increasing the range of the weather forecast
    
