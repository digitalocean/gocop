package main

import (
	"log"

	"github.com/digitalocean/gocop/action"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	action.Execute()
}
