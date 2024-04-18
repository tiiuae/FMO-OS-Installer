package screen

import (
	"ghaf-installer/global"

	"github.com/pterm/pterm"
)

// Method to get the heading message of screen
func (m ScreensMethods) CheckHardwareHeading() string {
	return "Check Hardware Compatibility"
}

func (m ScreensMethods) CheckHardware() {

	keyboardLayoutList := []string{"us", "fi"}
	// Select keyboard layout
	keyboardLayout, _ := pterm.DefaultInteractiveSelect.
		WithMaxHeight(20).
		WithOptions(keyboardLayoutList).
		Show("Select keyboard layout")
	_, err := global.ExecCommand(
		"sudo",
		"loadkeys",
		"-u",
		keyboardLayout,
	)

	if err == 0 {
		pterm.Info.Printfln("Load keyboard layout successfully")
	} else {
		pterm.Error.Printfln("Failed to load keyboard layout")
	}
	goToScreen(GetCurrentScreen() + 1)

}
