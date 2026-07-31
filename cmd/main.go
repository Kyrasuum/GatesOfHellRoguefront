package main

import (
	"log"

	"roguefront/internal/app"
)

func main() {
	a := app.NewApp()
	if a != nil {
		err := a.Start(false)
		if err != nil {
			log.Printf("ERR: %+v\n", err)
		}
	}
}
