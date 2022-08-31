package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/michaelcoll/gallery-daemon/indexer"
	"log"
)

var path = flag.String("f", ".", "The folder where the photos are.")

func main() {

	flag.Parse()

	fmt.Printf("Monitoring folder %s \n", color.GreenString(*path))

	err := indexer.Scan(*path)
	if err != nil {
		log.Fatal(err)
	}
}
