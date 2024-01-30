package screen

import (
	"ghaf-installer/global"
	"time"

	"github.com/pterm/pterm"
)

// Method to get the heading message of screen
func (m ScreensMethods) MountScreenHeading() string {
	return "Mount partition for configuring"
}

func (m ScreensMethods) MountScreen() {

	if !(haveInstalledSystem) {
		screenControlOption := appendScreenControl(make([]string, 0))
		// Print options to select device to install image
		selectedOption, _ := pterm.DefaultInteractiveSelect.
			WithOptions(screenControlOption).
			Show("No system installed, select what to do: ")
		if checkSkipScreen(selectedOption) {
			return
		}
		return
	}

	ghafMountingSpinner, _ := pterm.DefaultSpinner.
		WithShowTimer(false).
		WithRemoveWhenDone(true).
		Start("Mounting Partition")

	// Mount ghaf system
	mountGhaf("/dev/" + selectedPartition)

	// Wait time for user to read the message
	time.Sleep(2)
	ghafMountingSpinner.Stop()

	pterm.Info.Printfln("Ghaf has been mounted to /root")

	time.Sleep(1)
	goToScreen(GetCurrentScreen() + 1)
	return
}

func mountGhaf(disk string) {
	_, err := global.ExecCommand("mkdir", "-p", mountPoint)
	if err != 0 {
		panic(err)
	}

	_, err = global.ExecCommand("sudo", "mount", disk+"p2", mountPoint)
	if err != 0 {
		panic(err)
	}
	haveMountedSystem = true
}
