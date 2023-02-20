package main

import (
	"fmt"
	"log"

	"github.com/BillyBones007/loyalty-service/internal/app/client"
	"github.com/BillyBones007/loyalty-service/internal/app/server"
)

func main() {
	server := server.NewServer()
	client := client.NewAccrualClient(server.Storage, server.Config.AddrAccrual)
	go client.Run()
	fmt.Println("Configuration server:")
	fmt.Printf("Server is running on: %s\n", server.Config.AddrServ)
	fmt.Printf("DSN string: %s\n", server.Config.AddrDB)
	fmt.Printf("Accrual system address: %s\n", server.Config.AddrAccrual)
	defer server.ShutdownServer()
	log.Fatal(server.HTTPServer.ListenAndServe())
}
