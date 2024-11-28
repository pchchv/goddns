package main

import (
	"flag"

	"github.com/fatih/color"
	"github.com/pchchv/goddns/internal/utils"
)

var (
	optHelp = flag.Bool("h", false, "Show help")
	// Version is current version of GoDDNS.
	Version = "v0.1"
)

func main() {
	utils.Version = Version

	flag.Parse()
	if *optHelp {
		color.Cyan(utils.Logo, Version)
		flag.Usage()
		return
	}
}
