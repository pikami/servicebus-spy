package main

import (
	"flag"
	"log"

	"github.com/pikami/servicebus-spy/internal/amqphandler"
	messagecollection "github.com/pikami/servicebus-spy/internal/message_collection"
	servicebusclient "github.com/pikami/servicebus-spy/internal/servicebus_client"
	"github.com/pikami/servicebus-spy/internal/tcpproxy"
	webapi "github.com/pikami/servicebus-spy/internal/web_api"
)

func main() {
	serviceBusHost := flag.String("service-bus-host", "127.0.0.1:5672", "The host of the service bus")
	serviceBusConnectionString := flag.String("service-bus-connection-string", "", "The connection string to the service bus")
	proxyPort := flag.Int("proxy-port", 5666, "The port to listen for incoming connections")
	webPort := flag.Int("web-port", 8080, "The port to listen for incoming web requests")
	flag.Parse()

	if *serviceBusConnectionString == "" {
		log.Fatal("The service bus connection string is required")
	}

	messageCollection := messagecollection.NewMessageCollection()

	serviceBusClient := servicebusclient.NewServiceBusClient(*serviceBusConnectionString)
	defer serviceBusClient.Close()

	webAPI := webapi.NewWebAPI(*webPort, messageCollection, serviceBusClient)
	go webAPI.StartWebAPI()

	amqpHandler := amqphandler.NewAMQPHandler(messageCollection)
	p := tcpproxy.NewTcpProxy(*proxyPort, *serviceBusHost, amqpHandler)
	p.Start()
}
