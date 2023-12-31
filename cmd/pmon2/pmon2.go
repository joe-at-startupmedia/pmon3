package main

import (
	"log"
	"pmon2/cli/cmd"
	"pmon2/pmond"
	"pmon2/pmond/conf"
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
