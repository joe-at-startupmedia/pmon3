package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tabalt/gracehttp"
	"log"
	"os"
)

func main() {
	addr := fmt.Sprintf("0.0.0.0:%s", os.Getenv("TEST_APP_PORT"))
	err := gracehttp.ListenAndServe(addr, gin.Default())
	if err != nil {
		log.Fatal(err)
	}
}
