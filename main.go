package main

import (
	"awesomeProject/router"
	"log"
)

func main() {
	r := router.InitRouter()

	err := r.Run()
	if err != nil {
		log.Panic(err)
	}
}
