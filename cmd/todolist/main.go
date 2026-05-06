package main

import (
	"flag"

	"github.com/max-messenger/max-bot-example-todolist/internal/app"
)

func main() {
	var (
		cfgName string
	)

	flag.StringVar(&cfgName, "c", "config.yaml", "config")
	flag.Parse()

	app.CreateApp(cfgName).Run()
}
