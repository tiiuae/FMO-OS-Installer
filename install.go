package main

import (
	"ghaf-installer/screen"
	"reflect"
	"time"
	"log"
	"os"

	"github.com/pterm/pterm"
)

func showcase(title string, seconds int, content func()) {
	pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgBlue)).
		WithTextStyle(pterm.NewStyle(pterm.FgWhite, pterm.BgBlue, pterm.Bold)).
		WithFullWidth().
		Println(title)
	pterm.Println()
	time.Sleep(time.Second / 2)
	content()
	time.Sleep(time.Second * time.Duration(seconds))
	print("\033[H\033[2J")
}

func main() {
	// use for retrieving methods
	var methodlist screen.ScreensMethods

	// Open log file
	logFile, err := os.OpenFile("/tmp/installer.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Can't open log file: %v", err)
	}
	defer logFile.Close()

	// Setup logger
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime)

	// Look for all screens in folder ./screen
	screen.InitScreen()
	for (screen.GetCurrentScreen()) < len(screen.Screens) {
		currentScreen := screen.GetCurrentScreen()

		screenHeading := reflect.ValueOf(methodlist).
			MethodByName(screen.Screens[currentScreen] + "Heading").
			Call(nil)[0].
			Interface().(string)

		screenFunc := reflect.ValueOf(methodlist).
			MethodByName(screen.Screens[currentScreen]).
			Interface().(func())

		showcase(screenHeading, 2, screenFunc)
	}

}
