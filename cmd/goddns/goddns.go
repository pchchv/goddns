package main

import (
	"flag"
	"os"

	"github.com/fatih/color"
	"github.com/pchchv/goddns/internal/utils"
)

const configEnv = "CONFIG"

var (
	optConf = flag.String("c", "./config.json", "Specify a config file")
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

	configPath := *optConf
	// read config path from the environment
	if os.Getenv(configEnv) != "" {
		// overwrite the config path
		configPath = os.Getenv(configEnv)
	}
}
