package main

import (
	"pmon2/pmond"
	"pmon2/pmond/conf"
	"pmon2/cli/cmd"
	"log"
)


func main() {
	config := conf.GetDefaultConf()
	err := pmond.Instance(config)
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Exec()
	if err != nil {
		log.Fatal(err)
	}
}
