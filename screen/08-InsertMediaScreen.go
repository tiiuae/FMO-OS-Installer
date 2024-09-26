package screen

import (
	"ghaf-installer/global"
	"strings"
	"time"
	"os"

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
	pterm.Info.Printfln("selected option: %s", selectedOption)

	for selectedOption == updateDriversStr {
		// If a skip option selected
		if checkSkipScreen(selectedOption) {
			pterm.Info.Printfln("Skip containers preload...")
			time.Sleep(3)
			return
		}

		selectedOption = SelectOption()
		pterm.Info.Printfln("selected option: %s", selectedOption)
	}

	// If a skip option selected
	if checkSkipScreen(selectedOption) {
		pterm.Info.Printfln("Skip containers preload...")
		time.Sleep(3)
		return
	}

	/***************** resize rootfs *******************/
	resizeRootFs("/dev/" + selectedPartition, mountPoint)

	/***************** mount media *********************/
	pterm.Info.Printfln("Use %s for preloaded containers", selectedOption)
	dev := strings.TrimSpace(
		strings.Split(string(selectedOption), global.SPACE_CHAR)[0],
	)
	pterm.Info.Printfln("Use %s for preloaded containers", dev)

	ghafMountingSpinner, _ := pterm.DefaultSpinner.
		WithShowTimer(false).
		WithRemoveWhenDone(true).
		Start("Mounting Partition")

	// Mount media
	mountMedia("/dev/" + dev, "/media/fmoos-containers")

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
	copyData("/media/fmoos-containers/*", mountPoint + "/var/fogdata/preloaded/")

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
	pterm.Info.Printfln("mkdir -p %s", to)
	msg, err := global.ExecCommand("mkdir", "-p", to)
  for _, m := range msg {
		pterm.Info.Printfln(m)
	}
	if err != 0 {
		pterm.Info.Printfln("mkdir -p %s failed..", to)
		panic(err)
	}

	if false {
		pterm.Info.Printfln("sudo rsync -ah --progress %s %s", from, to)
		msg, err = global.ExecCommand("sudo", "rsync", "-ah", "--progress", from, to)
		for _, m := range msg {
			pterm.Info.Printfln(m)
		}
		if err != 0 {
			pterm.Info.Printfln("sudo rsync -ah --progress %s %s failed..", from, to)
			panic(err)
		}
	}

	if true {
		pterm.Info.Printfln("sudo sh -c 'cp -r %s %s'", from, to)
		msg, err = global.ExecCommand("sudo", "sh", "-c", "cp -r " + from + " " + to)
  	for _, m := range msg {
			pterm.Info.Printfln(m)
		}
		if err != 0 {
			pterm.Info.Printfln("sudo sh -c 'cp -r %s %s' failed..", from, to)
			panic(err)
		}
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

func resizeRootFs(disk string, mountPoint string) {
	var resizeCmd = "(echo d; echo 2; echo n; echo 2; echo ''; echo ''; echo w;) | sudo fdisk " + disk +
									"\nsudo e2fsck -y -f " + disk + "p2" +
									"\nsudo resize2fs -f " + disk + "p2"

	umountMedia(mountPoint)

	// Write the string to the file
  err := os.WriteFile("resize.sh", []byte(resizeCmd), 0777)
  if err != nil {
      pterm.Info.Printfln("Error writing to file:", err)
      return
  }

	pterm.Info.Printfln("sudo bash ./resize.sh")
  msg, er := global.ExecCommand("sudo", "bash", "./resize.sh")
  for _, m := range msg {
		pterm.Info.Printfln(m)
	}
	if er != 0 {
		pterm.Info.Printfln("sudo bash ./resize.sh failed..")
		panic(er)
	}

	if false {
		pterm.Info.Printfln("sudo e2fsck -f %s", disk + "p2")
		msg, er = global.ExecCommand("sudo", "e2fsck", "-f ", disk + "p2")
		for _, m := range msg {
			pterm.Info.Printfln(m)
		}
		if er != 0 {
			pterm.Info.Printfln("sudo e2fsck -f %s failed..", disk + "p2")
			panic(er)
		}

		pterm.Info.Printfln("sudo resize2fs %s", disk + "p2")
		msg, er = global.ExecCommand("sudo", "resize2fs ", disk + "p2")
		for _, m := range msg {
			pterm.Info.Printfln(m)
		}
		if er != 0 {
			pterm.Info.Printfln("sudo resize2fs %s failed..", disk + "p2")
			panic(er)
		}
	}

	mountMedia(disk + "p2", mountPoint)
}

func mountMedia(disk string, mountPoint string) {
	pterm.Info.Printfln("mkdir -p %s", mountPoint)
	_, err := global.ExecCommand("mkdir", "-p", mountPoint)
	if err != 0 {
		pterm.Info.Printfln("mkdir -p %s failed..", mountPoint)
		panic(err)
	}

	pterm.Info.Printfln("sudo mount %s %s", disk, mountPoint)
	_, err = global.ExecCommand("sudo", "mount", disk, mountPoint)
	if err != 0 {
	pterm.Info.Printfln("sudo mount %s %s failed..", disk, mountPoint)
		panic(err)
	}
}

func umountMedia(mountPoint string) {
	pterm.Info.Printfln("sudo umount %s", mountPoint)
	_, err := global.ExecCommand("sudo", "umount", mountPoint)
	if err != 0 {
		pterm.Info.Printfln("sudo umount %s failed", mountPoint)
		panic(err)
	}
}
