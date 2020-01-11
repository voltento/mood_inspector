package main

import "mood_inspector/internal/app"

func main() {
	application := app.NewApp()
	application.Run()
}
