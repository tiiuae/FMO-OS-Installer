package screen

import (
	"bufio"
	"ghaf-installer/global"
	"log"
	"os"
	"strings"

	"github.com/pterm/pterm"
)

// Method to get the heading message of screen
func (m ScreensMethods) ImageScreenHeading() string {
	return "Select image"
}

func (m ScreensMethods) ImageScreen() {
	// Image list is a string with format:

	imageList, imageMap := readOSSFile(string(global.OSSfile))

	// Add skip options and print option list to select which image to install
	imageList = appendScreenControl(imageList)
	selectedImage, _ := pterm.DefaultInteractiveSelect.
		WithOptions(imageList).
		Show("Please select image to install")

	// If skip option is selected
	if checkSkipScreen(selectedImage) {
		return
	}

	global.Image2Install = imageMap[selectedImage]

	goToScreen(GetCurrentScreen() + 1)

}

func readOSSFile(OSSPath string) ([]string, map[string]string) {

	file, err := os.Open(OSSPath)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	lines := make([]string, 0)
	imageNames := make([]string, 0)
	imageMaps := make(map[string]string)
	// Read through 'tokens' until an EOF is encountered.
	for sc.Scan() {
		lines = append(lines, sc.Text())
		imageName := strings.Split(sc.Text(), string(";"))[0]
		imagePath := strings.Split(sc.Text(), string(";"))[1]
		imageNames = append(imageNames, imageName)
		imageMaps[imageName] = imagePath
	}

	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}
	return imageNames, imageMaps
}
