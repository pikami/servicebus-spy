package messagecollection

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type MessageCollection struct {
	mu       sync.Mutex
	messages []Message
}

type ApplicationDataType string

const (
	ApplicationDataTypeString  ApplicationDataType = "string"
	ApplicationDataTypeJSON    ApplicationDataType = "json"
	ApplicationDataTypeBytes   ApplicationDataType = "bytes"
	ApplicationDataTypeUnknown ApplicationDataType = "unknown"
)

type Message struct {
	LinkInfo MessageLinkInfo `json:"linkInfo"`

	MessageAnnotations map[string]any `json:"messageAnnotations"`
	MessageProperties  map[string]any `json:"messageProperties"`

	ApplicationData       any                 `json:"applicationData"`
	ApplicationDataType   ApplicationDataType `json:"applicationDataType"`
	ApplicationProperties map[string]any      `json:"applicationProperties"`

	Raw string `json:"raw"`
}

type MessageLinkInfo struct {
	Role          string `json:"role"`
	SourceAddress string `json:"sourceAddress"`
	TargetAddress string `json:"targetAddress"`
}

func NewMessageCollection() *MessageCollection {
	return &MessageCollection{
		messages: make([]Message, 0),
	}
}

func (c *MessageCollection) AddMessage(message Message) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.messages = append(c.messages, message)
}

func (c *MessageCollection) GetMessages() []Message {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.messages
}

func (c *MessageCollection) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.messages = make([]Message, 0)
}

func (c *MessageCollection) DumpToFile() {
	c.mu.Lock()
	defer c.mu.Unlock()

	json, err := json.Marshal(c.messages)
	if err != nil {
		fmt.Printf("Failed to marshal messages: %+v\n", err)
		return
	}

	os.WriteFile("messages.json", json, 0644)
}

func (c *MessageCollection) LoadFromFile() {
	jsonData, err := os.ReadFile("messages.json")
	if err != nil {
		fmt.Printf("Failed to read messages.json: %+v\n", err)
		return
	}

	var messages []Message
	err = json.Unmarshal(jsonData, &messages)
	if err != nil {
		fmt.Printf("Failed to unmarshal messages: %+v\n", err)
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.messages = messages
}
