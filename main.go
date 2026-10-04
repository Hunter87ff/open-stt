package main

import (
	"open-stt/lib/config"
	"open-stt/lib/controller"
)

func main() {
	config.StartASRThread()
	controller.Listen()
}
