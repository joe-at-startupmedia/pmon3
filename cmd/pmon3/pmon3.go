package main

import (
	"log"
	"pmon3/cli/cmd"
	"pmon3/pmond"
	"pmon3/pmond/conf"
)

func main() {
	err := pmond.Instance(conf.GetConfigFile())
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Exec()
	if err != nil {
		log.Fatal(err)
	}
}
