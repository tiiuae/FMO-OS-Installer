package screen

import (
	"bufio"
	"bytes"
	"ghaf-installer/global"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/pterm/pterm"
)

// Method to get the heading message of screen
func (m ScreensMethods) InsertMediaScreenHeading() string {
	return "Select docker preloaded media to install"
}

func (m ScreensMethods) InsertMediaScreen() {
	var drivesList []string
	var drivesListHeading string

	drivesList = appendScreenControl(drivesList)

	// If no images are selected to install
	if len(global.Image2Install) == 0 {
		pterm.Error.Printfln("No image is selected for the installation")
	} else {
		// Get all block devices
		drives, _ := global.ExecCommand("lsblk", "-d", "-e7")
		if len(drives) > 0 {
			drivesListHeading = "  " + drives[0]
			for _, d := range drives[1:len(drives)-1] {
				if strings.Contains(d, "nvme") {
					drivesList = append(drivesList, d)
				}
			}
		}
	}

	// Print options to select device to install image
	selectedOption, _ := pterm.DefaultInteractiveSelect.
		WithOptions(drivesList).
		Show("Please select device to install Ghaf \n  " + drivesListHeading)

	// If a skip option selected
	if checkSkipScreen(selectedOption) {
		return
	}

	/***************** Start Installing *******************/

	goToScreen(GetCurrentScreen() + 1)
	return
}
