package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// Compress the backups in a zip file
func archive() error {
	af, err := createArchive()
	if err != nil {
		return err
	}
	defer af.Close()

	if err = recurArchive(af, bcfg.Destination, ""); err != nil {
		return err
	}

	err = os.RemoveAll(bcfg.Destination)
	if err != nil {
		return err
	}

	return nil
}

// Create a zip file in the backup destination
func createArchive() (*zip.Writer, error) {
	af, err := os.Create(fmt.Sprintf("%s.zip", bcfg.Destination))
	if err != nil {
		return nil, err
	}

	return zip.NewWriter(af), nil
}

// Recursively add directories/files to the zip file
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
				log.Println("Failed to open ", cPath)
				return err
			}

			dstPath := strings.TrimPrefix(fmt.Sprintf("%s/%s", zipBase, f.Name()), "/")

			dst, err := aw.Create(dstPath)
			if err != nil {
				return err
			}

			bs := bufio.NewScanner(src)

			if *encryption {
				// TODO: Change to env variables or config file
				sbh := genSBH("test", 42, 1729)
				err = encryptFile(bs.Bytes(), dst, sbh)
				if err != nil {
					return err
				}
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
