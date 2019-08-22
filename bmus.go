package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
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

	flags := flag.String("f", "", "Flags to pass to rsync")
	target := flag.String("t", "", "Target dir/file to backup")
	dest := flag.String("d", "", "Destination for backup(s)")
	flag.Parse()

	if *flags == "" {
		*flags = "-az"
	}
	if *target == "" {
		*target = fmt.Sprintf("%s", cuser.HomeDir)
	}
	if *dest == "" {
		*dest = fmt.Sprintf("%s/Documents/backups", cuser.HomeDir)
	} else {
		*dest = fmt.Sprintf("%s/backups", *dest)
	}

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

func checkBackupDest() error {
	_, err := os.Stat(bcfg.Destination)
	if err != nil {
		os.MkdirAll(bcfg.Destination, 0777)
	}

	return nil
}

// Back Me Up Scotty
func BMUS() error {
	err = exec.Command("rsync", bcfg.Flags, bcfg.Target, bcfg.Destination).Run()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if err = BMUS(); err != nil {
		log.Printf("Failed to BMUS: %v", err)
	}
}
