package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
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

	fmt.Println("Writing Thunderstore metadata files")
	tsm := btmd.CreateThunderstoreManifestObject(webUrl, repackager)
	tsm.WriteMetadataToFolder(repackager.BuildDir, btmd)

	fmt.Println()
	fmt.Println(fmt.Sprintf("Place an icon.png (256x256) in the %s folder", repackager.BuildDir))
	fmt.Println("Once you have done this, select all files in that folder and zip them")
	fmt.Println("(On Windows: Select all -> Right click -> Send to -> Compressed (zipped) folder")
	fmt.Println()
	fmt.Println("When zipped, you can upload to Thunderstore.")
	fmt.Println()
}

func verifyUrlFormat(url string) {
	matched, _ := regexp.MatchString("^https?://bonetome\\.com/.+\\d+/?$", url)
	if !matched {
		log.Fatal("URL is not in the specified format")
	}
}
