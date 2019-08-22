package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func archive() error {
	af, err := createArchive(bcfg.Destination)
	if err != nil {
		return err
	}
	defer af.Close()

	if err := recurArchive(af, bcfg.Destination, ""); err != nil {
		return err
	}

	err = os.RemoveAll(bcfg.Destination)
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

			dstPath := fmt.Sprintf("%s/%s", zipBase, f.Name())
			dst, err := aw.Create(strings.TrimPrefix(dstPath, "/"))
			if err != nil {
				return err
			}

			_, err = io.Copy(dst, src)
			if err != nil {
				return err
			}
		} else if f.IsDir() {
			newBase := fmt.Sprintf("%s/%s", base, f.Name())
			newZipBase := fmt.Sprintf("%s/%s", zipBase, f.Name())

			err = recurArchive(aw, newBase, newZipBase)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
