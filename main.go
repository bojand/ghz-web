package main

import "github.com/bojand/ghz-web/config"

func main() {
	config, err := config.Read("")
	if err != nil {
		panic(err)
	}

	app := Application{
		Config: config,
	}

	app.Start()
}
