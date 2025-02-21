package utils

import (
	"os"
	"strconv"
)

/*
*
* Port/0 attempts to get the port from the envireonment variable PORT
* if the varaible is not defined, it will default to the port 8080
* you can provide an environment variable in linux through "export PORT=9000". For other OSes, feel free to consult their documentation
*
 */
func Port() int {
	port := 8080
	portStr := os.Getenv("PORT")

	if n, e := strconv.Atoi(portStr); e == nil {
		port = n
	}
	return port
}
