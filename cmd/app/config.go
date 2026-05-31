package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	EnvPrefix = "SERVICEBUS_SPY_"
)

type AppConfig struct {
	ServiceBusHost             string
	ServiceBusConnectionString string
	ProxyPort                  int
	WebPort                    int
	WebDist                    string
}

func parseFlags() (config AppConfig) {
	serviceBusHost := flag.String("service-bus-host", "127.0.0.1:5672", "The host of the service bus")
	serviceBusConnectionString := flag.String("service-bus-connection-string", "", "The connection string to the service bus")
	proxyPort := flag.Int("proxy-port", 5666, "The port to listen for incoming connections")
	webPort := flag.Int("web-port", 8080, "The port to listen for incoming web requests")
	webDist := flag.String("web-dist", "", "Path to built web UI directory to serve")

	flag.Parse()
	err := setFlagsFromEnvironment()
	if err != nil {
		log.Fatalf("Failed to set flags from environment: %v", err)
	}

	if *serviceBusConnectionString == "" {
		log.Fatal("The service bus connection string is required")
	}

	if *webDist != "" {
		info, err := os.Stat(*webDist)
		if err != nil {
			log.Fatalf("Failed to access web-dist directory: %v", err)
		}
		if !info.IsDir() {
			log.Fatal("web-dist must be a directory")
		}
	}

	return AppConfig{
		ServiceBusHost:             *serviceBusHost,
		ServiceBusConnectionString: *serviceBusConnectionString,
		ProxyPort:                  *proxyPort,
		WebPort:                    *webPort,
		WebDist:                    *webDist,
	}
}

func setFlagsFromEnvironment() (err error) {
	flag.VisitAll(func(f *flag.Flag) {
		name := EnvPrefix + strings.ToUpper(strings.Replace(f.Name, "-", "_", -1))
		if value, ok := os.LookupEnv(name); ok {
			err2 := flag.Set(f.Name, value)
			if err2 != nil {
				err = fmt.Errorf("failed setting flag from environment: %w", err2)
			}
		}
	})

	return
}
