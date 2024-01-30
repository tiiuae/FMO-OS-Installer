package screen

import (
	"ghaf-installer/global"
	"time"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

func (m ScreensMethods) WelcomeScreen() {
	area, _ := pterm.DefaultArea.WithCenter().WithCenter().Start()
	for i := 0; i < 2; i++ {
		str, _ := pterm.DefaultBigText.WithLetters(
			putils.LettersFromStringWithStyle("FMO", pterm.FgLightGreen.ToStyle()),
			putils.LettersFromString("-OS")).Srender()
		area.Update(str)
		time.Sleep(time.Second)
	}
	goToScreen(GetCurrentScreen() + 1)
	return
}

func (m ScreensMethods) WelcomeScreenHeading() string {
	return global.WelcomeMsg
}
