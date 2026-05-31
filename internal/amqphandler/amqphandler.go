package amqphandler

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"

	messagecollection "github.com/pikami/servicebus-spy/internal/message_collection"
	"github.com/pikami/servicebus-spy/internal/prsr"
)

type AMQPHandler struct {
	mu                sync.RWMutex
	connections       map[string]*Connection
	messageCollection *messagecollection.MessageCollection
}

type Connection struct {
	mu       sync.RWMutex
	channels map[int16]*Channel
}

type Channel struct {
	mu    sync.RWMutex
	links map[uint]Link
}

type Link struct {
	sourceAddress string
	targetAddress string
	linkType      string
}

func NewAMQPHandler(messageCollection *messagecollection.MessageCollection) *AMQPHandler {
	return &AMQPHandler{
		connections:       make(map[string]*Connection),
		messageCollection: messageCollection,
	}
}

func (h *AMQPHandler) OnNewConnection(connID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.connections[connID] = &Connection{
		channels: make(map[int16]*Channel),
	}
}

func (h *AMQPHandler) OnConnectionClosed(connID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.connections, connID)
}

func (h *AMQPHandler) GetConnectionWriter(connID string) io.Writer {
	return &AMQPHandlerWriter{
		amqpHandler: h,
		connID:      connID,
	}
}

func (h *AMQPHandler) OnDataReceived(connID string, data []byte) {
	connection := h.getOrCreateConnection(connID)

	parsedFrame := prsr.ParseFrame(data)
	if !parsedFrame.Success {
		fmt.Printf("Failed to parse frame: %+v\n", parsedFrame)
		return
	}

	channel := connection.getOrCreateChannel(parsedFrame.Channel)

	for _, result := range parsedFrame.Results {
		resultMap, ok := result.(map[string]any)
		if !ok {
			continue
		}

		resultType, ok := resultMap["_type"]
		if !ok {
			continue
		}

		if _, ok := resultType.(prsr.DescriptorKey); !ok {
			continue
		}

		switch resultType {
		case prsr.DescriptorKeyATTACH:
			h.handleAttachFrame(channel, resultMap)
			return

		case prsr.DescriptorKeyTRANSFER:
			h.handleTransferFrame(channel, parsedFrame)
			return
		}
	}
}

func (h *AMQPHandler) handleAttachFrame(channel *Channel, resultMap map[string]any) {
	role, ok := resultMap["role"].(bool)
	if !ok {
		return
	}

	handle, ok := resultMap["handle"].(uint)
	if !ok {
		return
	}

	linkType := "receiver"
	if role {
		linkType = "sender"
	}

	sourceAddress, _ := getNestedString(resultMap, "source", "address")
	targetAddress, _ := getNestedString(resultMap, "target", "address")
	channel.createLink(handle, sourceAddress, targetAddress, linkType)
}

func (h *AMQPHandler) handleTransferFrame(channel *Channel, frame prsr.ParseFrameResult) {
	var transfer map[string]any
	var messageProperties map[string]any
	var messageAnnotations map[string]any
	var applicationProperties map[string]any
	var applicationData map[string]any

	for _, result := range frame.Results {
		resultMap, ok := result.(map[string]any)
		if !ok {
			continue
		}

		resultType, ok := resultMap["_type"]
		if !ok {
			continue
		}

		switch resultType {
		case prsr.DescriptorKeyTRANSFER:
			transfer = resultMap
		case prsr.DescriptorKeyMESSAGE_PROPERTIES:
			messageProperties = resultMap
		case prsr.DescriptorKeyMESSAGE_ANNOTATIONS:
			messageAnnotations = resultMap
		case prsr.DescriptorKeyAPPLICATION_PROPERTIES:
			applicationProperties = resultMap
		case prsr.DescriptorKeyAPPLICATION_DATA:
			applicationData = resultMap
		}
	}

	if transfer == nil || applicationData == nil {
		return
	}

	handle, ok := transfer["handle"].(uint)
	if !ok {
		return
	}

	link, ok := channel.links[handle]
	if !ok {
		fmt.Printf("Link not found: %d\n", handle)
		return
	}

	jsonFrame, err := json.Marshal(frame)
	if err != nil {
		fmt.Printf("Failed to marshal frame: %+v\n", err)
		return
	}

	parsedApplicationData := parseApplicationData(applicationData)
	applicationDataType := messagecollection.ApplicationDataTypeUnknown
	if _, ok := parsedApplicationData.(string); ok {
		applicationDataType = messagecollection.ApplicationDataTypeString
	} else if _, ok := parsedApplicationData.(map[string]any); ok {
		applicationDataType = messagecollection.ApplicationDataTypeJSON
	} else if _, ok := parsedApplicationData.([]byte); ok {
		applicationDataType = messagecollection.ApplicationDataTypeBytes
	}

	h.messageCollection.AddMessage(messagecollection.Message{
		LinkInfo: messagecollection.MessageLinkInfo{
			Role:          link.linkType,
			TargetAddress: link.targetAddress,
		},
		MessageAnnotations:    messageAnnotations,
		MessageProperties:     messageProperties,
		ApplicationData:       parsedApplicationData,
		ApplicationDataType:   applicationDataType,
		ApplicationProperties: applicationProperties,
		Raw:                   string(jsonFrame),
	})
}

func parseApplicationData(applicationData map[string]any) any {
	applicationDataBytes, ok := applicationData["_fields"].([]byte)
	if !ok {
		return []byte(applicationDataBytes)
	}

	if len(applicationDataBytes) == 0 {
		return []byte(applicationDataBytes)
	}

	if applicationDataBytes[0] == 0x22 { // starts with quote
		return string(applicationDataBytes)
	}

	if applicationDataBytes[0] == 0x7b { // starts with {
		var jsonData map[string]any
		err := json.Unmarshal(applicationDataBytes, &jsonData)
		if err != nil {
			return []byte(applicationDataBytes)
		}
		return jsonData
	}

	if applicationDataBytes[0] == 0x5b { // starts with [
		var jsonData []any
		err := json.Unmarshal(applicationDataBytes, &jsonData)
		if err != nil {
			return []byte(applicationDataBytes)
		}
		return jsonData
	}

	return []byte(applicationDataBytes)
}

func (h *AMQPHandler) getOrCreateConnection(connID string) *Connection {
	h.mu.Lock()
	defer h.mu.Unlock()

	connection, ok := h.connections[connID]
	if !ok {
		connection = &Connection{
			channels: make(map[int16]*Channel),
		}
		h.connections[connID] = connection
	}
	return connection
}

func (c *Connection) getOrCreateChannel(channelID int16) *Channel {
	c.mu.Lock()
	defer c.mu.Unlock()

	channel, ok := c.channels[channelID]
	if !ok {
		channel = &Channel{
			links: make(map[uint]Link),
		}
		c.channels[channelID] = channel
	}
	return channel
}

func (c *Channel) createLink(handle uint, sourceAddress string, targetAddress string, linkType string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.links[handle] = Link{
		sourceAddress: sourceAddress,
		targetAddress: targetAddress,
		linkType:      linkType,
	}
}

func getNestedString(m map[string]any, keys ...string) (string, bool) {
	var current any = m

	for _, key := range keys {
		nextMap, ok := current.(map[string]any)
		if !ok {
			return "", false
		}

		current, ok = nextMap[key]
		if !ok {
			return "", false
		}
	}

	s, ok := current.(string)
	return s, ok
}
