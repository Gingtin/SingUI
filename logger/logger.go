package logger

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type RingBuffer struct {
	mu      sync.RWMutex
	lines   []string
	maxSize int
	head    int
}

var (
	GlobalLogBuffer *RingBuffer
	clientsMu       sync.Mutex
	wsClients       = make(map[*websocket.Conn]bool)
)

func InitLogger(maxSize int) {
	if maxSize <= 0 {
		maxSize = 500
	}
	GlobalLogBuffer = &RingBuffer{
		lines:   make([]string, 0, maxSize),
		maxSize: maxSize,
	}
}

func (rb *RingBuffer) Append(line string) {
	if rb == nil {
		return
	}
	rb.mu.Lock()
	defer rb.mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	formattedLine := fmt.Sprintf("[%s] %s", timestamp, line)

	if len(rb.lines) < rb.maxSize {
		rb.lines = append(rb.lines, formattedLine)
	} else {
		rb.lines[rb.head] = formattedLine
		rb.head = (rb.head + 1) % rb.maxSize
	}

	BroadcastLog(formattedLine)
}

func (rb *RingBuffer) GetAll() []string {
	if rb == nil {
		return nil
	}
	rb.mu.RLock()
	defer rb.mu.RUnlock()

	res := make([]string, 0, len(rb.lines))
	n := len(rb.lines)
	if n < rb.maxSize {
		res = append(res, rb.lines...)
	} else {
		for i := 0; i < n; i++ {
			idx := (rb.head + i) % n
			res = append(res, rb.lines[idx])
		}
	}
	return res
}

func RegisterWSClient(conn *websocket.Conn) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	wsClients[conn] = true
}

func UnregisterWSClient(conn *websocket.Conn) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	delete(wsClients, conn)
	_ = conn.Close()
}

func BroadcastLog(msg string) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	for conn := range wsClients {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
			delete(wsClients, conn)
			_ = conn.Close()
		}
	}
}

func LogInfo(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	log.Println("[INFO] " + msg)
	if GlobalLogBuffer != nil {
		GlobalLogBuffer.Append("[INFO] " + msg)
	}
}

func LogError(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	log.Println("[ERROR] " + msg)
	if GlobalLogBuffer != nil {
		GlobalLogBuffer.Append("[ERROR] " + msg)
	}
}
