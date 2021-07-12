package main

import (
	"archive/zip"
	"fmt"
	"github.com/kjk/lzmadec"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type PullService struct{
	BoneTomeMetadata *BoneTomeMetadata
	TempDirectory string
}

func (ps *PullService) download(btmd *BoneTomeMetadata) {
	resp, err := http.Get(btmd.Asset.DownloadUrl)
	if err != nil {
		log.Fatal("Unable to establish connection to download file: "+btmd.Asset.DownloadUrl, err)
	}
	defer resp.Body.Close()
	file, _ := os.Create(btmd.Asset.Name)
	defer file.Close()
	io.Copy(file, resp.Body)
}

func (ps *PullService) Do(btmd *BoneTomeMetadata) {
	_, lstatErr := os.Lstat(btmd.Asset.Name + ".temp_dir")
	if !(os.IsNotExist(lstatErr)) {
		err := os.RemoveAll(btmd.Asset.Name + ".temp_dir")
		if err != nil {
			log.Fatal(err)
		}
	}
	ps.download(btmd)
	nameLower := strings.ToLower(btmd.Asset.Name)
	if strings.HasSuffix(nameLower, ".zip") {
		err := ps.unzip(btmd)
		if err != nil {
			log.Fatal(err)
		}
	} else if strings.HasSuffix(nameLower, ".7z") {
		ps.unzip7z(btmd)
	}
	ps.BoneTomeMetadata = btmd
	ps.TempDirectory = btmd.Asset.Name + ".temp_dir"
	err := os.Remove(btmd.Asset.Name)
	if err != nil {
		log.Fatal(err)
	}
}

func (ps *PullService) unzip(btmd *BoneTomeMetadata) error {
	baseDir := btmd.Asset.Name + ".temp_dir"

	r, err := zip.OpenReader(btmd.Asset.Name)
	if err != nil {
		return err
	}

	defer r.Close()

	for _, f := range r.File {

		fpath := filepath.Join(baseDir, f.Name)

		if !strings.HasPrefix(fpath, filepath.Clean(baseDir)+string(os.PathSeparator)) {
			return errors.New(fmt.Sprintf("%s is an invalid filepath", fpath))
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())

		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}

func (ps *PullService) unzip7z(btmd *BoneTomeMetadata) {
	archive, archErr := lzmadec.NewArchive(path.Join(".", btmd.Asset.Name))
	if archErr != nil {
		log.Fatal(archErr)
	}
	baseDir := path.Join(".", btmd.Asset.Name+".temp_dir")
	_ = os.Mkdir(baseDir, 0775)
	for entry := range archive.Entries {
		_ = archive.ExtractToFile(path.Join(baseDir, archive.Entries[entry].Path), archive.Entries[entry].Path)
	}
}
