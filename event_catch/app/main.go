package main

import (
	"event/app/app"
	"event/config"
	"flag"
	"fmt"
)

var configFlag = flag.String("config", "./config.toml", "toml env file not found")

func main() {
	flag.Parse()

	config.NewConfig(*configFlag)

	a := app.NewApp(config.NewConfig(*configFlag))
	fmt.Println(a)
}
