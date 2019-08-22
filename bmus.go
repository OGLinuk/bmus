package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"os/user"
)

var (
	configName = "config.json"
)

type BMUSConfig struct {
	Flags       string
	Target      string
	Destination string
}

var (
	cfg   BMUSConfig
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

	if *flags == "" && *target == "" && *dest == "" {
		cfg, err = loadConfig()
		if err != nil {
			log.Printf("Failed to LoadConfig: %v", err)
		}
	} else {
		if *flags == "" {
			*flags = "-azvhP"
		}
		if *target == "" {
			*target = fmt.Sprintf("%s", cuser.HomeDir)
		}
		if *dest == "" {
			*dest = fmt.Sprintf("%s/Documents/backups", cuser.HomeDir)
		}

		cfg = BMUSConfig{
			Flags:       *flags,
			Target:      *target,
			Destination: fmt.Sprintf("%s/backups", *dest),
		}

		if err = checkBackupDest(); err != nil {
			log.Printf("Failed to checkDefaultBackup: %v", err)
		}
	}
}

func main() {
	if err = BMUS(); err != nil {
		log.Printf("Failed to BMUS: %v", err)
	}
}

// Back Me Up Scotty
func BMUS() error {
	err = exec.Command("rsync", cfg.Flags, cfg.Target, cfg.Destination).Run()
	if err != nil {
		return err
	}

	return nil
}
