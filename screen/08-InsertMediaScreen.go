package screen

import (
	"ghaf-installer/global"
	"strings"
	"time"

	"github.com/pterm/pterm"
)

var updateDriversStr = "Update drivers list"

// Method to get the heading message of screen
func (m ScreensMethods) InsertMediaScreenHeading() string {
	return "Select docker preloaded media to install\n"
}

func (m ScreensMethods) InsertMediaScreen() {
	/***************** check installaion ***************/
	if !(haveInstalledSystem) || !(haveMountedSystem) {
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

	/***************** select media ********************/
	selectedOption := SelectOption()
	pterm.info.printfln("selected option: %s\n", selectedOption)

	for selectedOption == updateDriversStr {
		// If a skip option selected
		if checkSkipScreen(selectedOption) {
			pterm.Info.Printfln("Skip containers preload...")
			time.Sleep(3)
			return
		}

		selectedOption = SelectOption()
		pterm.info.printfln("selected option: %s\n", selectedOption)
	}

	/***************** select media ********************/
	/***************** mount media *********************/
	ghafMountingSpinner, _ := pterm.DefaultSpinner.
		WithShowTimer(false).
		WithRemoveWhenDone(true).
		Start("Mounting Partition")

	// Mount ghaf system
	mountMedia("/dev/" + selectedPartition, "/media/fmoos-containers")

	// Wait time for user to read the message
	time.Sleep(2)
	ghafMountingSpinner.Stop()

	pterm.Info.Printfln("Containers has been mounted..")
	time.Sleep(1)

	/***************** start copying *******************/
	ghafMountingSpinner, _ = pterm.DefaultSpinner.
		WithShowTimer(false).
		WithRemoveWhenDone(true).
		Start("Mounting Partition")

	// Umount media
	copyData("/media/fmoos-containers" + "/*", mountPoint + "/var/fogdata/preloaded")

	// Wait time for user to read the message
	time.Sleep(2)
	ghafMountingSpinner.Stop()

	pterm.Info.Printfln("Containers has been copied..")
	time.Sleep(1)

	/***************** umount media ********************/
	ghafMountingSpinner, _ = pterm.DefaultSpinner.
		WithShowTimer(false).
		WithRemoveWhenDone(true).
		Start("Mounting Partition")

	// Umount media
	umountMedia("/media/fmoos-containers")

	// Wait time for user to read the message
	time.Sleep(2)
	ghafMountingSpinner.Stop()

	pterm.Info.Printfln("Containers has been umounted..")
	time.Sleep(1)

	goToScreen(GetCurrentScreen() + 1)
	return
}

func copyData(from string, to string) {
	_, err := global.ExecCommand("mkdir", "-p", to)
	if err != 0 {
		panic(err)
	}

	_, err = global.ExecCommand("sudo", "cp", "-r", from, to)
	if err != 0 {
		panic(err)
	}
}

func SelectOption() string {
	var drivesList []string
	var drivesListHeading string

	drivesList = appendScreenControl(drivesList)
	drivesList = append(drivesList, updateDriversStr)

	// Get all block devices
	drives, _ := global.ExecCommand("lsblk", "-d", "-e7", "-o", "name,label")
	if len(drives) > 0 {
		drivesListHeading = "  " + drives[0]
		for _, d := range drives[1:len(drives)-1] {
			if strings.Contains(d, "fmoos-containers") {
				drivesList = append(drivesList, d)
			}
		}
	}

	// Print options to select device to install image
	selectedOption, _ := pterm.DefaultInteractiveSelect.
		WithOptions(drivesList).
		Show("Please select device with FMO-OS preloaded containers\n" + drivesListHeading)

	return selectedOption
}

func mountMedia(disk string, mounPoint string) {
	_, err := global.ExecCommand("mkdir", "-p", mountPoint)
	if err != 0 {
		panic(err)
	}

	_, err = global.ExecCommand("sudo", "mount", disk, mountPoint)
	if err != 0 {
		panic(err)
	}
}

func umountMedia(mounPoint string) {
	_, err := global.ExecCommand("sudo", "umount", mountPoint)
	if err != 0 {
		panic(err)
	}
}
