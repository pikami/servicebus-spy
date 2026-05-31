package amqphandler

type AMQPHandlerWriter struct {
	amqpHandler *AMQPHandler
	connID      string
}

func (w *AMQPHandlerWriter) Write(p []byte) (n int, err error) {
	w.amqpHandler.OnDataReceived(w.connID, p)
	return len(p), nil
}
