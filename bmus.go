package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
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
			*dest = fmt.Sprintf("%s/Documents", cuser.HomeDir)
		}

		cfg = BMUSConfig{
			Flags:       *flags,
			Target:      *target,
			Destination: *dest,
		}
	}
}

func main() {
	if err := BMUS(); err != nil {
		log.Printf("Failed to BMUS: %v", err)
	}
}

// Back Me Up Scotty
func BMUS() error {

	err = exec.Command("rsync", cfg.Flags, cfg.Target, cfg.Destination).Run()
	if err != nil {
		return nil
	}

	return nil
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

func loadConfig() (BMUSConfig, error) {
	var bcfg BMUSConfig

	f, err := os.Open(configName)
	if err != nil {
		saveConfig(&BMUSConfig{
			Flags:       "-azvhP",
			Target:      cuser.HomeDir,
			Destination: fmt.Sprintf("%s/Documents", cuser.HomeDir),
		})
		return cfg, err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	err = decoder.Decode(&bcfg)

	return bcfg, err
}
