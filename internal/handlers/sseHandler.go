package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
	"tiny-drop/internal/utils"
)

type Client struct {
	ip          string
	fingerprint string
	ch          chan string
}

var clientsByIP = make(map[string]map[string]*Client)
var clientsMu sync.RWMutex

func SSEHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		utils.SendError(w, http.StatusInternalServerError, "Streaming Unsupported!", nil)
		return
	}

	ip := utils.GetUserIp(r)

	query := r.URL.Query()
	fingerprint := query.Get("fingerprint")
	// key := fmt.Sprintf("%s|%s", ip, fingerprint)

	ch := make(chan string)

	client := &Client{
		ip:          ip,
		fingerprint: fingerprint,
		ch:          ch,
	}

	clientsMu.Lock()
	if clientsByIP[ip] == nil {
		clientsByIP[ip] = make(map[string]*Client)
	}
	clientsByIP[ip][fingerprint] = client
	clientsMu.Unlock()

	log.Println("CLient connected: ",  ip, fingerprint)

	defer func() {
		clientsMu.Lock()
		delete(clientsByIP[ip], fingerprint)
		if len(clientsByIP[ip]) == 0 {
			delete(clientsByIP, ip)
		}
		clientsMu.Unlock()
		log.Println("Client disconnected:",  ip, fingerprint)
	}()

	// channel for when client disconnects
	notify := r.Context().Done()

	for {
		select {
		case <-notify:
			return
		case msg := <-ch:
			_, _ = fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		}
	}
}

func BroadcastToIP(ip string, message string) {
	clientsMu.RLock()
	defer clientsMu.RUnlock()

	for _, client := range clientsByIP[ip] {
		select {
		case client.ch <- message:
		default:
			// Prevent blocking if client is slow or dead
		}
	}
}

func SSEHandlerTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		utils.SendError(w, http.StatusInternalServerError, "Streaming Unsupported!", nil)
		return
	}

	// channel for when client disconnects
	notify := r.Context().Done()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-notify:
			{
				log.Println("Client disconnected")
				return
			}
		// send event
		case <-ticker.C:
			{
				eventData := time.Now().Format(time.RFC3339)
				_, err := w.Write([]byte("data: " + eventData + "\n\n"))
				if err != nil {
					log.Printf("Error writing SSE: %v", err)
					return
				}
				flusher.Flush()
			}
		}
	}
}
