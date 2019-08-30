package main

import (
	"os"
	"os/exec"
)

// Ensures that the backup destination exists; if not then create one
func checkBackupDest() error {
	_, err := os.Stat(bcfg.Destination)
	if err != nil {
		os.MkdirAll(bcfg.Destination, 0777)
	}

	return nil
}

// Rsync wrapper that executes on given (t)arget
func scotty(t string) error {
	err = exec.Command("rsync", bcfg.Flags, t, bcfg.Destination).Run()
	if err != nil {
		return err
	}

	return nil
}
