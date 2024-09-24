package screen

import (
	"ghaf-installer/global"
	"strings"

	"github.com/pterm/pterm"
)

var updateDriversStr = "Update drivers list"

// Method to get the heading message of screen
func (m ScreensMethods) InsertMediaScreenHeading() string {
	return "Select docker preloaded media to install\n"
}

func (m ScreensMethods) InsertMediaScreen() {
	selectedOption := SelectOption()

	for selectedOption != updateDriversStr {
		// If a skip option selected
		if checkSkipScreen(selectedOption) {
			return
		}

		selectedOption = SelectOption()
	}



	/***************** Start Installing *******************/

	goToScreen(GetCurrentScreen() + 1)
	return
}

func string SelectOption() {
	var drivesList []string
	var drivesListHeading string

	drivesList = appendScreenControl(drivesList)

	// Get all block devices
	drives, _ := global.ExecCommand("lsblk", "-d", "-e7", "-o", "name,label")
	if len(drives) > 0 {
		drivesListHeading = "  " + drives[0]
		for _, d := range drives[1:len(drives)-1] {
			if strings.Contains(d, "fmoos-containers") {
				drivesList = append(drivesList, d)
			}
		}
	} else {
		drivesList = append(drivesList, updateDriversStr)
	}

	// Print options to select device to install image
	selectedOption, _ := pterm.DefaultInteractiveSelect.
		WithOptions(drivesList).
		Show("Please select device with FMO-OS preloaded containers\n" + drivesListHeading)

	return selectedOption
}
