package main

import (
	"log"
	"pmon3/cli"
	"pmon3/cli/cmd"
	"pmon3/conf"
)

func main() {
	err := cli.Instance(conf.GetConfigFile())
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Exec()
	if err != nil {
		log.Fatal(err)
	}
}
