package main

import (
	"pmon3/pmond"
	"pmon3/pmond/conf"
	"pmon3/pmond/god"
)

func main() {
	err := pmond.Instance(conf.GetConfigFile())
	if err != nil {
		pmond.Log.Fatal(err)
	}
	god.New()
}
