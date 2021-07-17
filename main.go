package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	url := ""
	// Get and validate package URL
	fmt.Print("Enter the BoneTome page URL here (https://bonetome.com/h3vr/mods/id): ")
	reader := bufio.NewReader(os.Stdin)
	line, _, _ := reader.ReadLine()
	url = string(line)
	verifyUrlFormat(url)

	// Populate as much metadata as possible from BoneTome page.
	btmd := &BoneTomeMetadata{}
	btmd.PopulateFromUrl(url)

	// Add a description field to be used in the manifest.
	fmt.Println()
	fmt.Print("Please enter a short description of the package (250 characters or fewer): ")
	descReader := bufio.NewReader(os.Stdin)
	descLine, _, _ := descReader.ReadLine()
	btmd.ProvideDescription(string(descLine))

	// Get website_url field to be used in the manifest.
	fmt.Println()
	fmt.Print("Please enter a URL to be displayed on Thunderstore. Leave blank if none. This URL is usually a link to a GitHub repository: ")
	webUrlReader := bufio.NewReader(os.Stdin)
	webUrlLine, _, _ := webUrlReader.ReadLine()
	webUrl := string(webUrlLine)

	fmt.Println()
	fmt.Println(fmt.Sprintf("Downloading %s from BoneTome", btmd.Asset.Name))
	pullService := &PullService{}
	pullService.Do(btmd)

	fmt.Println("Repacking files")
	repackager := &Repackager{}
	repackager.Do(pullService)

	if repackager.ContainsDeliFiles && !repackager.ContainsLvoFiles {
		repackager.ContainsLvoFiles = askBooleanQuestion("Do you require OtherLoader?")
	}

	repackager.RequiresH3VRUtilities = askBooleanQuestion("Do you require H3VRUtilities?")

	fmt.Println("Writing Thunderstore metadata files")
	tsm := btmd.CreateThunderstoreManifestObject(webUrl, repackager)
	tsm.WriteMetadataToFolder(repackager.BuildDir, btmd)

	correspondingNumber := inputValidNumber()

	fmt.Println()

	switch correspondingNumber {
	case 1:
		fmt.Println(fmt.Sprintf("Place an icon.png (256x256) in the %s folder", repackager.BuildDir))
		fmt.Println("Once you have done this, select all files in that folder and zip them")
		fmt.Println("(On Windows: Select all -> Right click -> Send to -> Compressed (zipped) folder")
		fmt.Println()
		fmt.Println("When zipped, you can upload to Thunderstore.")
	case 2:
		zipName := fmt.Sprintf("%s-%s.zip", btmd.PackageName, btmd.Version)
		repackager.ZipBuildContents(zipName)
		fmt.Println(fmt.Sprintf("The zip file [%s] can be imported into a supported Thunderstore mod manager", zipName))
	}

	fmt.Println()
}

func inputValidNumber() int {
	fmt.Println()
	fmt.Println("Enter the corresponding number:")
	fmt.Println("(1). I'm going to upload my mod to Thunderstore.")
	fmt.Println("(2). I want to import a mod using a mod manager.")
	fmt.Print("Enter number: ")

	numReader := bufio.NewReader(os.Stdin)
	numLine, _ := numReader.ReadByte()
	switch string(numLine) {
	case "1":
		return 1
	case "2":
		return 2
	default:
		fmt.Println("The number inputted was invalid. Please select a correct number.")
		return inputValidNumber()
	}
}

func askBooleanQuestion(question string) bool {
	fmt.Println()
	fmt.Print(fmt.Sprintf("%s (y/n): ", question))

	numReader := bufio.NewReader(os.Stdin)
	numLine, _ := numReader.ReadByte()
	switch strings.ToLower(string(numLine)) {
	case "y":
		return true
	case "n":
		return false
	default:
		fmt.Println("The input was invalid. Please try again.")
		return askBooleanQuestion(question)
	}
}

func verifyUrlFormat(url string) {
	matched, _ := regexp.MatchString("^https?://bonetome\\.com/.+\\d+/?$", url)
	if !matched {
		log.Fatal("URL is not in the specified format")
	}
}
