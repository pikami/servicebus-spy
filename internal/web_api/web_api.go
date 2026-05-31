package webapi

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	messagecollection "github.com/pikami/servicebus-spy/internal/message_collection"
	servicebusclient "github.com/pikami/servicebus-spy/internal/servicebus_client"
)

type WebAPI struct {
	port              int
	messageCollection *messagecollection.MessageCollection
	serviceBusClient  *servicebusclient.ServiceBusClient
}

func NewWebAPI(port int, messageCollection *messagecollection.MessageCollection, serviceBusClient *servicebusclient.ServiceBusClient) *WebAPI {
	return &WebAPI{
		port:              port,
		messageCollection: messageCollection,
		serviceBusClient:  serviceBusClient,
	}
}

func (w *WebAPI) StartWebAPI() {
	router := gin.Default()
	router.GET("/messages", w.getMessages)
	router.GET("/messages/dump", w.dumpMessages)
	router.GET("/messages/load", w.loadMessages)
	router.GET("/messages/clear", w.clearMessages)
	router.POST("/messages/send", w.sendMessage)
	http.ListenAndServe(fmt.Sprintf(":%d", w.port), router)
}

func (w *WebAPI) getMessages(c *gin.Context) {
	messages := w.messageCollection.GetMessages()
	c.JSON(http.StatusOK, messages)
}

func (w *WebAPI) dumpMessages(c *gin.Context) {
	w.messageCollection.DumpToFile()
	c.JSON(http.StatusOK, "Messages dumped to file")
}

func (w *WebAPI) loadMessages(c *gin.Context) {
	w.messageCollection.LoadFromFile()
	c.JSON(http.StatusOK, "Messages loaded from file")
}

func (w *WebAPI) clearMessages(c *gin.Context) {
	w.messageCollection.Clear()
	c.JSON(http.StatusOK, "Messages cleared")
}

func (w *WebAPI) sendMessage(c *gin.Context) {
	var message servicebusclient.ServiceBusMessage
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := w.serviceBusClient.SendMessage(message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Message sent")
}
