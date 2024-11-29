package main

import (
	"flag"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/pchchv/goddns/internal/settings"
	"github.com/pchchv/goddns/internal/utils"
)

const configEnv = "CONFIG"

var (
	// Version is current version of GoDDNS.
	Version = "v0.1"
	config  settings.Settings
	optHelp = flag.Bool("h", false, "Show help")
	optConf = flag.String("c", "./config.json", "Specify a config file")
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

	// Load settings from configs file
	if err := settings.LoadSettings(configPath, &config); err != nil {
		log.Fatal(err)
	}

	if err := utils.CheckSettings(&config); err != nil {
		log.Fatal("Invalid settings: ", err.Error())
	}
}
