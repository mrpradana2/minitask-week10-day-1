package main

import (
	"log"
	"os"
	"tikcitz-app/internals/routes"
	"tikcitz-app/pkg"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	if err := pkg.Connect(); err != nil {
		log.Printf("unable to create connection pool:  %v\n", err)
		os.Exit(1)
	}

	// closing DB
	defer func()  {
		log.Println("Closing DB...")
		pkg.DB.Close()	
	}()

    router := routes.InitRouter()

    router.Run()
}