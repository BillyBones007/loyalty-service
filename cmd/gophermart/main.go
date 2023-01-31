package main

import (
	"fmt"
	"log"

	"github.com/BillyBones007/loyalty-service/internal/app/server"
)

func main() {
	app := server.NewServer()
	fmt.Println("Configuration server:")
	fmt.Printf("Server is running on: %s\n", app.Config.AddrServ)
	fmt.Printf("DSN string: %s\n", app.Config.AddrDB)
	fmt.Printf("Accrual system address: %s\n", app.Config.AddrAccrual)
	defer app.ShutdownServer()
	log.Fatal(app.HTTPServer.ListenAndServe())
}
