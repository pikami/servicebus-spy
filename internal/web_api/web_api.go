package webapi

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	messagecollection "github.com/pikami/servicebus-spy/internal/message_collection"
	servicebusclient "github.com/pikami/servicebus-spy/internal/servicebus_client"
)

type WebAPI struct {
	port              int
	webDist           string
	messageCollection *messagecollection.MessageCollection
	serviceBusClient  *servicebusclient.ServiceBusClient
}

func NewWebAPI(port int, webDist string, messageCollection *messagecollection.MessageCollection, serviceBusClient *servicebusclient.ServiceBusClient) *WebAPI {
	return &WebAPI{
		port:              port,
		webDist:           webDist,
		messageCollection: messageCollection,
		serviceBusClient:  serviceBusClient,
	}
}

func (w *WebAPI) StartWebAPI() {
	router := gin.Default()
	w.registerAPIRoutes(router)
	w.registerAPIRoutes(router.Group("/api"))

	if w.webDist != "" {
		w.serveStaticFiles(router)
		log.Printf("Serving web UI from %s on port %d", w.webDist, w.port)
	}

	http.ListenAndServe(fmt.Sprintf(":%d", w.port), router)
}

func (w *WebAPI) registerAPIRoutes(r gin.IRoutes) {
	r.GET("/messages", w.getMessages)
	r.GET("/messages/dump", w.dumpMessages)
	r.GET("/messages/load", w.loadMessages)
	r.GET("/messages/clear", w.clearMessages)
	r.POST("/messages/send", w.sendMessage)
}

func (w *WebAPI) serveStaticFiles(router *gin.Engine) {
	dir := http.Dir(w.webDist)
	fileServer := http.FileServer(dir)

	router.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}

		cleanPath := strings.TrimPrefix(c.Request.URL.Path, "/")
		f, err := dir.Open(cleanPath)
		if err == nil {
			defer f.Close()
			stat, statErr := f.Stat()
			if statErr == nil && !stat.IsDir() {
				fileServer.ServeHTTP(c.Writer, c.Request)
				return
			}
		}

		c.File(filepath.Join(w.webDist, "index.html"))
	})
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
