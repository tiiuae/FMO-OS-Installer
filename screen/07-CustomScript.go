package screen

import (
	"errors"
	"ghaf-installer/global"
	"os"
	"strings"
	"time"

	"github.com/pterm/pterm"
)

// Path to registration-agent-laptop binary
var customScript = "./script/app"
var customScriptHeading = "Custom Script"

// Environment variables required for registration-agent-laptop
var folderPaths = "/var/fogdata/certs;/var/fogdata/config;/var/fogdata/token"

// Method to get the heading message of screen
func (m ScreensMethods) CustomScriptHeading() string {
	return customScriptHeading
}

func (m ScreensMethods) CustomScript() {

	script_err := false
	command := strings.Split(customScript, global.SPACE_CHAR)[0]
	if _, err := os.Stat(command); errors.Is(err, os.ErrNotExist) {
		pterm.Error.Printfln("Custom script not found!")
		script_err = true
	}

	if !(haveMountedSystem) {
		pterm.Error.Printfln("No system mounted!")
		script_err = true
	}

	if script_err {
		screenControlOption := appendScreenControl(make([]string, 0))
		// Print options to select device to install image
		selectedOption, _ := pterm.DefaultInteractiveSelect.
			WithOptions(screenControlOption).
			Show("Select what to do: ")
		if checkSkipScreen(selectedOption) {
			return
		}
		return
	}

	// Set create folders to store certificates
	prepareEnvironment()

	// Execute registration-agent-laptop binary
	global.ExecCommandWithLiveMessage("bash", strings.Split(customScript, global.SPACE_CHAR)...)

	// Wait for 3 seconds for user to read the finish log
	time.Sleep(3)
	goToScreen(GetCurrentScreen() + 1)

}

func prepareEnvironment() {
	// Create folder for certificates and config

	paths := strings.Split(string(folderPaths), ";")
	for _, folderPath := range paths {
		_, err := global.ExecCommand("sudo", "mkdir", "-p", folderPath)
		if err != 0 {
			panic(err)
		}
	}
}
