package main

import (
	"log"
	"pmon3/pmond"
	"pmon3/pmond/conf"
	"pmon3/pmond/god"
)

func main() {
	err := pmond.Instance(conf.GetConfigFile())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("pmon3 daemon is running! \n")

	god.NewMonitor()
}
