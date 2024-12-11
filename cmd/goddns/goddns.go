package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/pchchv/goddns/internal/manager"
	"github.com/pchchv/goddns/internal/settings"
	"github.com/pchchv/goddns/internal/utils"
)

const configEnv = "CONFIG"

var (
	config  settings.Settings
	Version = "v0.1" // current version of GoDDNS
	optHelp = flag.Bool("h", false, "Show help")
	optConf = flag.String("c", "./config.json", "Specify a config file")
	optAddr = flag.String("a", ":9000", "Specify the address to listen on")
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

	// Create DNS manager
	dnsManager := manager.GetDNSManager(configPath, &config, *optAddr)

	// Run DNS manager
	log.Println("GoDDNS started, starting the DNS manager...")
	dnsManager.Run()

	// handle the signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// stop the DNS manager
	<-c
	log.Println("GoDDNS is terminated, stopping the DNS manager...")
	dnsManager.Stop()

	// wait for the goroutines to exit
	time.Sleep(200 * time.Millisecond)
	log.Println("GoDDNS is stopped, bye!")
}
