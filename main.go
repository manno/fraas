package main

import (
	"log"

	"manno.name/mm/faas/actions"
)

func main() {
	app := actions.App()
	if err := app.Serve(); err != nil {
		log.Fatal(err)
	}
}
