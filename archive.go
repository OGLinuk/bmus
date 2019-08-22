package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func archive() error {
	afname := fmt.Sprintf("%s/backup_%d", bcfg.Destination, time.Now().Unix())

	_, tname := filepath.Split(bcfg.Target)

	af, err := createArchive(afname)
	if err != nil {
		return err
	}
	defer af.Close()

	if err := recurArchive(af, bcfg.Destination, ""); err != nil {
		return err
	}

	err = os.RemoveAll(fmt.Sprintf("%s/%s", bcfg.Destination, tname))
	if err != nil {
		return err
	}

	return nil
}

func createArchive(archiveName string) (*zip.Writer, error) {
	af, err := os.Create(fmt.Sprintf("%s.zip", archiveName))
	if err != nil {
		return nil, err
	}

	return zip.NewWriter(af), nil
}

func recurArchive(aw *zip.Writer, base, zipBase string) error {
	files, err := ioutil.ReadDir(base)
	if err != nil {
		return err
	}

	for _, f := range files {
		cPath := fmt.Sprintf("%s/%s", base, f.Name())

		if !f.IsDir() {
			src, err := os.Open(cPath)
			if err != nil {
				return err
			}

			dst, err := aw.Create(fmt.Sprintf("%s/%s", zipBase, f.Name()))
			if err != nil {
				return err
			}

			_, err = io.Copy(dst, src)
			if err != nil {
				return err
			}
		} else if f.IsDir() {
			newBase := fmt.Sprintf("%s/%s", base, f.Name())

			err = recurArchive(aw, newBase, fmt.Sprintf("%s/%s", zipBase, f.Name()))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
