package screen

import (
	"ghaf-installer/global"
	"os"
	"strings"

	"github.com/pterm/pterm"
)

// Method to get the heading message of screen
func (m ScreensMethods) NetworkInterfacesHeading() string {
	return "Connect to network"
}

func (m ScreensMethods) NetworkInterfaces() {
	// List all network devices
	networkInterfaces, _ := os.ReadDir("/sys/class/net")
	var networkInterfaceList []string
	for _, file := range networkInterfaces {
		networkInterfaceList = append(networkInterfaceList, file.Name())
	}

	// Select network device
	networkInterface, _ = pterm.DefaultInteractiveSelect.
		WithMaxHeight(20).
		WithOptions(networkInterfaceList).
		Show("Select network interface")

	// If the device is wireless
	if _, err := os.Stat("/sys/class/net/" + networkInterface + "/wireless"); !os.IsNotExist(err) {
		WifiScreen()
		return
	}

	goToScreen(GetCurrentScreen() + 1)
}

func WifiScreen() {

	wifiConnectSpinner, _ := pterm.DefaultSpinner.
		WithShowTimer(false).
		WithRemoveWhenDone(true).
		Start("Scanning Wifi")

	wifiConnect, _ := global.ExecCommand(
		"nmcli",
		"-t",
		"--fields",
		"SSID,SIGNAL,SECURITY",
		"dev",
		"wifi",
	)
	wifiConnectSpinner.Stop()

	skipWifi := false
	if len(wifiConnect) == 0 {
		skipWifi = true
		goToScreen(GetCurrentScreen() + 1)
		return
	}
	for !ConnectionStatus && !skipWifi {
		wifiList := wifiConnect[0 : len(wifiConnect)-1]
		longestWifiSSID := 0
		for _, wifi := range wifiList {
			wifiMessage := strings.Split(strings.Replace(string(wifi), "\\:", ":", -1), ":")
			wSSID := strings.Join(wifiMessage[0:len(wifiMessage)-2], ":")
			if len(wSSID) > longestWifiSSID {
				longestWifiSSID = len(wSSID)
			}
		}
		var wifiListBeautified []string

		for _, wifi := range wifiList {
			wifiMessage := strings.Split(strings.Replace(string(wifi), "\\:", ":", -1), ":")
			wSSID := strings.Join(wifiMessage[0:len(wifiMessage)-2], ":")
			wSignal := wifiMessage[len(wifiMessage)-2]
			wSecurity := wifiMessage[len(wifiMessage)-1]
			wSSIDBeautified := wSSID + strings.Repeat(
				" ",
				longestWifiSSID+2-len(wSSID),
			)
			wSignalBeautified := wSignal + strings.Repeat(" ", 8-len(wSignal))
			wifiListBeautified = append(
				wifiListBeautified,
				wSSIDBeautified+"||"+wSignalBeautified+"||"+wSecurity,
			)
		}

		wifiListBeautified = appendScreenControl(wifiListBeautified)
		wifiListHeading := "SSID" + strings.Repeat(
			" ",
			longestWifiSSID+2-len("SSID"),
		) + "||SIGNAL  ||SECURITY"
		selectedWifi, _ := pterm.DefaultInteractiveSelect.
			WithMaxHeight(20).
			WithOptions(wifiListBeautified).
			Show("Wifi list \n  " + wifiListHeading)

		if checkSkipScreen(selectedWifi) {
			skipWifi = true
			return
		}

		SSID := strings.TrimSpace(
			strings.Split(string(selectedWifi), string("||"))[0],
		)
		pterm.Info.Printfln("Connect to %s", SSID)

		password, _ := pterm.DefaultInteractiveTextInput.
			WithMultiLine(false).
			WithMask("*").
			Show("Password (If no password, leave empty)")

		wifiConnectingSpinner, _ := pterm.DefaultSpinner.
			WithShowTimer(false).
			WithRemoveWhenDone(true).
			Start("Connecting to " + SSID)
		_, err := global.ExecCommand(
			"nmcli",
			"dev",
			"wifi",
			"connect",
			SSID,
			"password",
			password,
		)
		wifiConnectingSpinner.Stop()
		if err == 0 {
			ConnectionStatus = true
			pterm.Info.Printfln("Connected")
		} else {
			ConnectionStatus = false
			pterm.Error.Printfln("Failed to connect to " + SSID)
		}
	}

	goToScreen(GetCurrentScreen() + 1)

}
