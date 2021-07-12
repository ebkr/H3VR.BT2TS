package main

import (
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

type Repackager struct {
	BuildDir                string
	ContainsLvoFiles        bool
	ContainsDeliFiles       bool
	ContainsSideloaderFiles bool
}

func (rpkg *Repackager) Do(ps *PullService) {
	baseDir := ps.BoneTomeMetadata.Asset.Name + "._build"

	_, lstatErr := os.Lstat(baseDir)
	if !(os.IsNotExist(lstatErr)) {
		err := os.RemoveAll(baseDir)
		if err != nil {
			log.Fatal(err)
		}
	}

	err := os.Mkdir(baseDir, 0775)
	if err != nil {
		log.Fatal(err)
	}

	rpkg.BuildDir = baseDir

	files := rpkg.getFiles(ps.TempDirectory)
	rpkg.moveFiles(files)

	err = os.RemoveAll(ps.TempDirectory)
	if err != nil {
		log.Fatal(err)
	}
}

func (rpkg *Repackager) moveFiles(files []string) {

	fileNameRegexp := regexp.MustCompilePOSIX("[^/\\\\]+$")

	lvoPath := path.Join("plugins", "LegacyVirtualObjects")
	deliPath := path.Join("plugins", "DeliMods")
	fileLocations := map[string][]string{
		"Sideloader": {},
		"plugins":    {},
		lvoPath:      {},
		deliPath:     {},
	}
	for fileIndex := range files {
		fileLower := strings.ToLower(files[fileIndex])
		if strings.HasSuffix(fileLower, ".hotmod") || strings.HasSuffix(fileLower, ".h3mod") {
			fileLocations["Sideloader"] = append(fileLocations["Sideloader"], files[fileIndex])
			rpkg.ContainsSideloaderFiles = true
		} else if strings.HasSuffix(fileLower, ".dll") {
			fileLocations["plugins"] = append(fileLocations["plugins"], files[fileIndex])
		} else if strings.HasSuffix(fileLower, ".deli") {
			fileLocations[deliPath] = append(fileLocations[deliPath], files[fileIndex])
			rpkg.ContainsDeliFiles = true
		} else {
			fileName := fileNameRegexp.FindString(files[fileIndex])
			if !strings.Contains(fileName, ".") {
				// LVO
				fileLocations[lvoPath] = append(fileLocations[lvoPath], files[fileIndex])
				rpkg.ContainsLvoFiles = true
			} else {
				if strings.HasSuffix(fileLower, ".manifest") {
					nameWithoutExtension := fileName[:len(fileName)-len(".manifest")]
					fileDir := files[fileIndex][:len(files[fileIndex])-len(fileName)]

					_, osErr := os.Lstat(path.Join(fileDir, nameWithoutExtension))
					if !os.IsNotExist(osErr) {
						fileLocations[lvoPath] = append(fileLocations[lvoPath], files[fileIndex])
					}
				} else {
					fileLocations["plugins"] = append(fileLocations["plugins"], files[fileIndex])
				}
			}
		}
	}
	for k, v := range fileLocations {
		if len(v) == 0 {
			continue
		}
		err := os.MkdirAll(path.Join(rpkg.BuildDir, k), 0775)
		if err != nil {
			log.Fatal(err)
		}
		for fileIndex := range v {
			// Get file name, stripping path and separators.
			fileName := fileNameRegexp.FindString(v[fileIndex])
			file, err := os.Create(path.Join(rpkg.BuildDir, k, fileName))
			defer file.Close()
			if err != nil {
				log.Fatal(err)
			}
			reader, err := os.Open(v[fileIndex])
			defer reader.Close()
			if err != nil {
				log.Fatal(err)
			}
			_, err = io.Copy(file, reader)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

// https://stackoverflow.com/a/49196644
func (rpkg *Repackager) getFiles(dir string) []string {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return files
}
