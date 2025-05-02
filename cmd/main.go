package main

import (
	"log"
	"os"
	"tikcitz-app/internals/routes"
	"tikcitz-app/pkg"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	db, err := pkg.Connect()
	if  err != nil {
		log.Printf("unable to create connection pool:  %v\n", err)
		os.Exit(1)
	}

	// closing DB
	defer func()  {
		log.Println("Closing DB...")
		db.Close()	
	}()
	
	// var hash pkg.HashConfig
	// hash.UseDefaultConfig()
	// password := "kucinggarong"
	// hashedPassword, _ := hash.GenHashedPassword(password)
	// log.Println("[DEBUG] password: ", password)
	// log.Println("[DEBUG] hash: ", hashedPassword)
	rdb := pkg.RedisConnect()
    router := routes.InitRouter(db, rdb)
	router.Static("/img", "./public/img")

    router.Run()
}