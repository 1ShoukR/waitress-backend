// This is the entry point of the Application
//
// It creates a new server instance and starts the server


package main

import (
	"fmt"
	"waitress-backend/internal/server"
)

func main() {
	serverInstance := server.NewServer() // Renamed to avoid shadowing the package name

	err := serverInstance.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
