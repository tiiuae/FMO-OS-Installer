package screen

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pterm/pterm"
)

type ScreensMethods struct{}

// Defines constants and variables for installation process
const nextScreenMsg = ">>--Skip to next step------->>"
const exitInstallerMsg = ">>--Quit installation------->>"
const previousScreenMsg = "<<--Back to previous step---<<"

var ConnectionStatus = false
var selectedPartition string
var haveInstalledSystem = false
var haveMountedSystem = false

var currentInstallationScreen = 0
var Screens = make(map[int]string)

var sourceDir = "."
var screenDir = sourceDir + "/screen"
var dataDir = sourceDir + "/data"
var mountPoint = "/home/huy/root"

var envPath = "./.env"
var networkInterface = ""

// Docker URL selections
var dockerURLs = "ghcr.io*cr.airoplatform.com"
var dockerURLPath = "/var/fogdata/cr.url"

func GetCurrentScreen() int {
	return currentInstallationScreen
}

func goToScreen(screen int) {
	currentInstallationScreen = screen
}

func appendScreenControl(messages []string) []string {
	messages = append(
		[]string{exitInstallerMsg},
		messages...)
	messages = append(
		[]string{nextScreenMsg},
		messages...)
	messages = append(
		[]string{previousScreenMsg},
		messages...)
	return messages
}

func checkSkipScreen(input string) bool {
	if input == nextScreenMsg {
		goToScreen(GetCurrentScreen() + 1)
		return true
	} else if input == previousScreenMsg {
		goToScreen(GetCurrentScreen() - 1)
		return true
	} else if input == exitInstallerMsg {
		// Go to Exit Screen
		goToScreen(len(Screens)-1)
		return true
	}
	return false
}

func InitScreen() {
	files, err := os.ReadDir(screenDir)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		// Split file name into arrays/slices
		// "00-WelcomeScreen.go" => ["00" "WelcomeScreen"]
		fileName := strings.Split(
			file.Name()[:len(file.Name())-len(filepath.Ext(file.Name()))],
			"-",
		)

		// Skip if file name not following above format (e.g "screen.go")
		if len(fileName) != 2 {
			continue
		}

		order, _ := strconv.Atoi(fileName[0])
		Screens[order] = fileName[1]
	}
}

func RefreshScreen(title string) {

	print("\033[H\033[2J")

	pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgBlue)).
		WithTextStyle(pterm.NewStyle(pterm.FgWhite, pterm.BgBlue, pterm.Bold)).
		WithFullWidth().
		Println(title)
	pterm.Println()

}
