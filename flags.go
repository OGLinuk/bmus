package main

import (
	"flag"
	"fmt"
)

var (
	flags      = flag.String("f", "", "Flags to pass to rsync")
	target     = flag.String("t", "", "Target dir/file to backup")
	dest       = flag.String("d", "", "Destination for backup(s)")
	encryption = flag.Bool("enc", false, "Encrypt backup file")
)

func parseFlags() {
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
}
