package app

import (
	"log"
	"pilulia_bot/logger/consts"
)

func InitApp() {
	app, err := NewApp()
	if err != nil {
		log.Fatal(consts.ErrorApp, err)
	}
	app.Start()
}
