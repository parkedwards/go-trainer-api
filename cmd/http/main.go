package main

import (
	app "github.com/parkedwards/go-trainer-api/pkg"
)

func main() {
	r := app.Init()
	app.Boot(r)
}
