package main

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
)

type ThunderstoreVersion struct {
	FullName string `json:"full_name"`
}

type ThunderstorePackage struct {
	FullName string `json:"full_name"`
	Versions []ThunderstoreVersion `json:"versions"`
}

type ThunderstoreApi struct {
	Packages []ThunderstorePackage
}

func (tsapi *ThunderstoreApi) Pull() {
	resp, err := http.Get("https://h3vr.thunderstore.io/api/v1/package/")
	if err != nil {
		log.Fatal("Unable to establish a connection to the TS API:",  err)
	}
	defer resp.Body.Close()

	var packages []ThunderstorePackage
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	json.Unmarshal(body, &packages)
	tsapi.Packages = packages
}

func (tsapi *ThunderstoreApi) GetDependencyStringFor(dep string) (string, error) {
	for pkgIndex := range tsapi.Packages {
		if tsapi.Packages[pkgIndex].FullName == dep {
			// V1 API returns packages in order of upload. Almost always guaranteed to be latest.
			return tsapi.Packages[pkgIndex].Versions[0].FullName, nil
		}
	}
	return "", errors.New(fmt.Sprintf("unable to find dependency string for %s", dep))
}