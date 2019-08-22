package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func loadConfig() (BMUSConfig, error) {
	var bcfg BMUSConfig

	f, err := os.Open(configName)
	if err != nil {
		saveConfig(&BMUSConfig{
			Flags:       "-azvhP",
			Target:      cuser.HomeDir,
			Destination: fmt.Sprintf("%s/Documents/backups", cuser.HomeDir),
		})
		return cfg, err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	err = decoder.Decode(&bcfg)

	return bcfg, err
}

func saveConfig(bcfg *BMUSConfig) error {
	f, err := os.Create(configName)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "\t")
	encoder.Encode(bcfg)

	return nil
}

func checkBackupDest() error {
	log.Printf("Checking if %s exists ...", cfg.Destination)

	_, err := os.Stat(cfg.Destination)
	if err != nil {
		os.MkdirAll(cfg.Destination, 0777)
	}

	return nil
}
