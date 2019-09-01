package main

import (
	"log"
	"os/user"
	"path/filepath"
)

type BMUSConfig struct {
	Flags       string
	Target      string
	Destination string
}

var (
	bcfg  BMUSConfig
	cuser *user.User
	err   error
)

func init() {
	cuser, err = user.Current()
	if err != nil {
		log.Printf("Failed to get user.Current: %v", err)
	}

	parseFlags()
	targetAbsPath, err := filepath.Abs(*target)
	if err != nil {
		log.Printf("Failed to get target absolute path ...")
	}

	destAbsPath, err := filepath.Abs(*dest)
	if err != nil {
		log.Printf("Failed to get dest absolute path ...")
	}

	bcfg = BMUSConfig{
		Flags:       *flags,
		Target:      targetAbsPath,
		Destination: destAbsPath,
	}

	if err = checkBackupDest(); err != nil {
		log.Printf("Failed to checkDefaultBackup: %v", err)
	}
}

// Back Me Up Scotty
func BMUS() error {
	if err = scotty(bcfg.Target); err != nil {
		return err
	}

	if err = archive(); err != nil {
		return err
	}

	return nil
}

func main() {
	if err = BMUS(); err != nil {
		log.Printf("Failed to BMUS: %v", err)
	}
}
