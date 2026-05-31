package servicebusclient

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

type ServiceBusMessage struct {
	QueueOrTopic string
	Body         string
	ContentType  string
	Subject      string
}

type ServiceBusClient struct {
	client *azservicebus.Client
}

func NewServiceBusClient(connectionString string) *ServiceBusClient {
	client, err := azservicebus.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		log.Fatalf("Failed to create service bus client: %v", err)
	}
	return &ServiceBusClient{client: client}
}

func (c *ServiceBusClient) Close() {
	c.client.Close(context.Background())
}

func (c *ServiceBusClient) SendMessage(message ServiceBusMessage) error {
	sender, err := c.client.NewSender(message.QueueOrTopic, nil)
	if err != nil {
		return err
	}
	defer sender.Close(context.Background())

	return sender.SendMessage(context.Background(), &azservicebus.Message{
		Body:        []byte(message.Body),
		ContentType: &message.ContentType,
		Subject:     &message.Subject,
	}, nil)
}
