package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
)

type BoneTomeAsset struct {
	Name        string
	DownloadUrl string
}

type BoneTomeMetadata struct {
	RealName    string
	PackageName string
	Version     string
	Asset       *BoneTomeAsset
	Readme      string
	Description string
}

type ThunderstoreManifest struct {
	Name         string   `json:"name"`
	Version      string   `json:"version_number"`
	Description  string   `json:"description"`
	WebsiteUrl   string   `json:"website_url"`
	Dependencies []string `json:"dependencies"`
}

func (btmd *BoneTomeMetadata) PopulateFromUrl(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Connection could not be made to the BoneTome server.")
	}
	defer resp.Body.Close()

	doc, docErr := goquery.NewDocumentFromReader(resp.Body)
	if docErr != nil {
		log.Fatal(docErr)
	}

	btmd.RealName = strings.TrimSpace(doc.Find("#info-assets-top").Find("[class='assets-title']").Text())
	btmd.PackageName = strings.ReplaceAll(btmd.RealName, " ", "_")
	unsafeReg := regexp.MustCompilePOSIX("[^a-zA-Z0-9_]")
	btmd.PackageName = unsafeReg.ReplaceAllString(btmd.PackageName, "_")
	underscoreReg := regexp.MustCompilePOSIX("_+")
	btmd.PackageName = underscoreReg.ReplaceAllString(btmd.PackageName, "_")
	btmd.PackageName = strings.Trim(btmd.PackageName, "_")

	fmt.Println()
	fmt.Println("Package name:", btmd.PackageName)

	bta := BoneTomeAsset{}
	btmd.Asset = &bta

	modInfoChildrenContainer := doc.Find("#mod-data").Children()
	for i := range modInfoChildrenContainer.Nodes {
		node := modInfoChildrenContainer.Eq(i)
		switch strings.ToLower(node.Find("[class='mod-data-header']").Text()) {
		case "version":
			btmd.Version = strings.TrimSpace(node.Find("[class='mod-data-content']").Text())
			break
		case "file name":
			bta.Name = strings.TrimSpace(node.Find("[class='mod-data-content']").Text())
			break
		}
	}

	bta.DownloadUrl, _ = doc.Find("#download-button").Parent().Attr("href")
	bta.DownloadUrl = "https://bonetome.com" + bta.DownloadUrl

	btmd.buildReadme(doc)

}

/** Not a perfect generation. Doesn't capture headings.
BT page HTML isn't the best and has loose inner text coupled with potential elements.
Similar to:
	<div>
		text
		<span>more text</span>
		even more text
	</div>
*/
func (btmd *BoneTomeMetadata) buildReadme(doc *goquery.Document) {
	btmd.Readme = "# " + btmd.RealName + "\n\n"
	btmd.Readme = btmd.Readme + doc.Find("#mod-info-desc").Eq(0).Text()
	changelog := doc.Find("#mod-info-desc").Eq(1)
	if strings.ToLower(changelog.Text()) != "n/a" {
		btmd.Readme = btmd.Readme + "\n\n## Changelog" + "\n\n"
		btmd.Readme = btmd.Readme + doc.Find("#mod-info-desc").Eq(1).Text()
	}
}

func (btmd *BoneTomeMetadata) ProvideDescription(desc string) {
	btmd.Description = desc
}

func (btmd *BoneTomeMetadata) CreateThunderstoreManifestObject(websiteUrl string, rpkg *Repackager) *ThunderstoreManifest {
	tsm := ThunderstoreManifest{
		Name:         btmd.PackageName,
		Version:      btmd.Version,
		Description:  btmd.Description,
		WebsiteUrl:   websiteUrl,
		Dependencies: make([]string, 0),
	}

	tsapi := &ThunderstoreApi{}
	tsapi.Pull()

	if rpkg.ContainsLvoFiles {
		dep, err := tsapi.GetDependencyStringFor("devyndamonster-OtherLoader")
		if err != nil {
			log.Fatal(err)
		}
		tsm.Dependencies = append(tsm.Dependencies, dep)
	}
	if rpkg.ContainsSideloaderFiles {
		dep, err := tsapi.GetDependencyStringFor("denikson-H3VR_Sideloader")
		if err != nil {
			log.Fatal(err)
		}
		tsm.Dependencies = append(tsm.Dependencies, dep)
	}
	if rpkg.ContainsDeliFiles {
		dep, err := tsapi.GetDependencyStringFor("DeliCollective-Deli")
		if err != nil {
			log.Fatal(err)
		}
		tsm.Dependencies = append(tsm.Dependencies, dep)
	}

	return &tsm
}

func (tsm *ThunderstoreManifest) WriteMetadataToFolder(dir string, btmd *BoneTomeMetadata) {

	marshalledData, marshalErr := json.Marshal(&tsm)
	if marshalErr != nil {
		log.Fatal(marshalErr)
	}

	var out bytes.Buffer
	jsonErr := json.Indent(&out, marshalledData, "", "	")
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	manifestFile, err := os.Create(path.Join(dir, "manifest.json"))
	if err != nil {
		log.Fatal(err)
	}
	defer manifestFile.Close()

	_, copyErr := io.Copy(manifestFile, bytes.NewReader(out.Bytes()))
	if copyErr != nil {
		log.Fatal(copyErr)
	}

	readmeFile, err := os.Create(path.Join(dir, "README.md"))
	if err != nil {
		log.Fatal(err)
	}
	defer readmeFile.Close()

	_, copyErr = io.Copy(readmeFile, bytes.NewReader([]byte(btmd.Readme)))
	if copyErr != nil {
		log.Fatal(copyErr)
	}

}
