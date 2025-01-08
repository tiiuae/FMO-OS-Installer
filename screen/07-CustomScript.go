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
	logWithPriority("DEBUG", "Prepare installation environment. In..")

	paths := strings.Split(string(folderPaths), ";")
	logWithPriority("DEBUG", "paths: %v", paths)

	checkEnvFile := strings.Split(paths[0], "/")
	logWithPriority("DEBUG", "checkEnvFile: %v", checkEnvFile)

	if checkEnvFile[len(checkEnvFile)-1] == ".env" {
		logWithPriority("DEBUG", ".env path has been found")

		envPath = strings.Split(string(folderPaths), ";")[0]
		logWithPriority("DEBUG", "envPath: %v", envPath)

		f, err := os.OpenFile(envPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			logWithPriority("ERROR", "Can not open file: %v", envPath)
			panic(err)
		}

		defer f.Close()

		if _, err = f.WriteString("NETWORK_INTERFACE=" + networkInterface); err != nil {
			panic(err)
		}
	} else {
		logWithPriority("ERROR", ".env path has been found")
	}

	// Symlink to simulate installed system
	symlink := paths[1]
	// Source location of installed system
	source := paths[2]

	logWithPriority("DEBUG", "Create new path: %v", mountPoint+source)
	_, err := global.ExecCommand("sudo", "mkdir", "-p", mountPoint+source)
	if err != 0 {
		logWithPriority("ERROR", "Can not create dir: %v",  mountPoint+source)
		panic(err)
	}

	logWithPriority("DEBUG", "Create new symlink: %v, %v", mountPoint+source, symlink)
	_, err = global.ExecCommand("sudo", "ln", "-s", mountPoint+source, symlink)
	if err != 0 {
		logWithPriority("ERROR", "Can not create symlink: %v, %v", mountPoint+source, symlink)
		panic(err)
	}

	// Create child folders for certificates and config
	logWithPriority("DEBUG", "Create child folders for certificates and config")
	paths = paths[3:]
	for _, folderPath := range paths {
		logWithPriority("DEBUG", "Create new path: %v", mountPoint+folderPath)
		_, err := global.ExecCommand("sudo", "mkdir", "-p", mountPoint+folderPath)
		if err != 0 {
			logWithPriority("ERROR", "Can not create dir: %v",  mountPoint+folderPath)
			panic(err)
		}
	}
	logWithPriority("DEBUG", "Prepare installation environment. Out..")
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
