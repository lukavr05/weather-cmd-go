This is my command line program that provides relevant weather information in terminal using Go and the OpenWeather API

Providing a city in the UK as an argument will retreive the information for that city. Without an argument, the default is set to London

This program also uses flags to specify whether a forecast or extended weather report is desired
    
    -f              specifies that a weather report is being requested, defaulting to the current day's forecast
    -e              specifies an extended weather report, supplying additional information to the current weather report as well as increasing the range of the weather forecast
    -days           an int specifying the number of days to display for the extended weather forecast (between 1 and 5)
    -default        displays the current default city stored in the config file
    -setdefault     updates the default value in the config file

This program makes use of a configuration file to store the default city that is used when no specific city is given as an argument. Without setting the default or having the config file in the path, the system will use "London" as its default. This can be changed in the program itself in the main function denoted by the system_default variable
