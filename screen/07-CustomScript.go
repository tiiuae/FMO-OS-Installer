package screen

import (
	"errors"
	"ghaf-installer/global"
	"os"
	"path/filepath"
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

	selectDockerURL()
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

	RefreshScreen(customScriptHeading)
	// Set create folders to store certificates
	prepareEnvironment()

	// Execute registration-agent-laptop binary
	global.ExecCommandWithLiveMessage("bash", strings.Split(customScript, global.SPACE_CHAR)...)

	// Wait for 3 seconds for user to read the finish log
	time.Sleep(3)
	goToScreen(GetCurrentScreen() + 1)

}

func prepareEnvironment() {
	paths := strings.Split(string(folderPaths), ";")
	checkEnvFile := strings.Split(paths[0], "/")
	if checkEnvFile[len(checkEnvFile)-1] == ".env" {
		envPath = strings.Split(string(folderPaths), ";")[0]
		f, err := os.OpenFile(envPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}

		defer f.Close()

		if _, err = f.WriteString("NETWORK_INTERFACE=" + networkInterface); err != nil {
			panic(err)
		}

		paths = paths[1:]
	}

	// Create folder for certificates and config
	for _, folderPath := range paths {
		_, err := global.ExecCommand("sudo", "mkdir", "-p", folderPath)
		if err != 0 {
			panic(err)
		}
	}
}

func selectDockerURL() {
	selectURL, _ := pterm.DefaultInteractiveSelect.
		WithOptions(strings.Split(dockerURLs, "*")).
		Show("Please select docker URL")

	_, err_int := global.ExecCommand("mkdir", "-p", filepath.Dir(mountPoint+dockerURLPath))
	if err_int != 0 {
		panic(err_int)
	}

	f, err := os.Create(mountPoint + dockerURLPath)
	if err != nil {
		panic(err)
	}

	defer func() {
		f.Close()
	}()

	_, err = f.WriteString(selectURL)
	if err != nil {
		panic(err)
	}

	err = f.Close()
	if err != nil {
		panic(err)
	}

	_, err_int = global.ExecCommand(
		"sudo",
		"chmod",
		"644",
		mountPoint+dockerURLPath,
	)
	if err_int != 0 {
		panic(err)
	}

}
