package main

import (
	"github.com/pikami/servicebus-spy/internal/amqphandler"
	messagecollection "github.com/pikami/servicebus-spy/internal/message_collection"
	servicebusclient "github.com/pikami/servicebus-spy/internal/servicebus_client"
	"github.com/pikami/servicebus-spy/internal/tcpproxy"
	webapi "github.com/pikami/servicebus-spy/internal/web_api"
)

func main() {
	config := parseFlags()

	messageCollection := messagecollection.NewMessageCollection()

	serviceBusClient := servicebusclient.NewServiceBusClient(config.ServiceBusConnectionString)
	defer serviceBusClient.Close()

	webAPI := webapi.NewWebAPI(config.WebPort, messageCollection, serviceBusClient)
	go webAPI.StartWebAPI()

	amqpHandler := amqphandler.NewAMQPHandler(messageCollection)
	p := tcpproxy.NewTcpProxy(config.ProxyPort, config.ServiceBusHost, amqpHandler)
	p.Start()
}
