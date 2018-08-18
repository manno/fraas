package main

import (
	"log"

	"manno.name/mm/fraas/actions"
)

func main() {
	app := actions.App()
	if err := app.Serve(); err != nil {
		log.Fatal(err)
	}
}
