package main

import (
	"pmon2/pmond"
	"pmon2/pmond/conf"
	"pmon2/pmond/god"
	"log"
)

func main() {
	config := conf.GetDefaultConf()
	err := pmond.Instance(config)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("pmon2 daemon is running! \n")

	god.NewMonitor()
}
