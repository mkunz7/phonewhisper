package main

import (
    "log"
    "net/http"
    "sync"
    "encoding/json"

    "github.com/gorilla/websocket"
)

type Connection struct {
    *websocket.Conn
    mu sync.Mutex
}

type Server struct {
    broadcaster     *Connection
    watchers       map[*Connection]bool
    waitingWatchers map[*Connection]bool
    mu             sync.RWMutex
}

func NewServer() *Server {
    return &Server{
        watchers:        make(map[*Connection]bool),
        waitingWatchers: make(map[*Connection]bool),
    }
}

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // Allow all origins for demo
    },
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("Websocket upgrade failed: %v", err)
        return
    }

    c := &Connection{Conn: conn}
    
    // Handle messages
    for {
        messageType, message, err := conn.ReadMessage()
        if err != nil {
            s.handleDisconnect(c)
            break
        }

        if messageType != websocket.TextMessage {
            continue
        }

        s.handleMessage(c, string(message))
    }
}

func (s *Server) handleMessage(conn *Connection, message string) {
    s.mu.Lock()
    defer s.mu.Unlock()

    var msg map[string]interface{}
    if err := json.Unmarshal([]byte(message), &msg); err != nil {
        log.Printf("Error parsing message: %v", err)
        return
    }

    messageType, ok := msg["type"].(string)
    if !ok {
        log.Printf("Message type not found")
        return
    }

    switch messageType {
    case "broadcaster":
        if s.broadcaster != nil {
            s.broadcaster.mu.Lock()
            s.broadcaster.WriteJSON(map[string]string{"type": "broadcasterDisconnected"})
            s.broadcaster.mu.Unlock()
        }
        s.broadcaster = conn
        
        // Notify all clients of new broadcaster
        for watcher := range s.watchers {
            watcher.mu.Lock()
            watcher.WriteJSON(map[string]string{"type": "broadcasterConnected"})
            watcher.mu.Unlock()
        }

        // Connect waiting watchers
        for watcher := range s.waitingWatchers {
            conn.mu.Lock()
            conn.WriteJSON(map[string]interface{}{
                "type": "watcher",
                "id":   watcher.Conn.RemoteAddr().String(),
            })
            conn.mu.Unlock()
            delete(s.waitingWatchers, watcher)
            s.watchers[watcher] = true
        }

    case "watcher":
        if s.broadcaster != nil {
            s.watchers[conn] = true
            s.broadcaster.mu.Lock()
            s.broadcaster.WriteJSON(map[string]interface{}{
                "type": "watcher",
                "id":   conn.Conn.RemoteAddr().String(),
            })
            s.broadcaster.mu.Unlock()
        } else {
            s.waitingWatchers[conn] = true
            conn.mu.Lock()
            conn.WriteJSON(map[string]string{"type": "noBroadcaster"})
            conn.mu.Unlock()
        }

    case "offer":
        if targetID, ok := msg["id"].(string); ok {
            // Forward the offer to the specific watcher
            for watcher := range s.watchers {
                if watcher.Conn.RemoteAddr().String() == targetID {
                    watcher.mu.Lock()
                    watcher.WriteJSON(map[string]interface{}{
                        "type": "offer",
                        "id":   conn.Conn.RemoteAddr().String(),
                        "description": msg["description"],
                    })
                    watcher.mu.Unlock()
                    break
                }
            }
        }

    case "answer":
        if targetID, ok := msg["id"].(string); ok {
            if s.broadcaster != nil && s.broadcaster.Conn.RemoteAddr().String() == targetID {
                s.broadcaster.mu.Lock()
                s.broadcaster.WriteJSON(map[string]interface{}{
                    "type": "answer",
                    "id":   conn.Conn.RemoteAddr().String(),
                    "description": msg["description"],
                })
                s.broadcaster.mu.Unlock()
            }
        }

    case "candidate":
        if targetID, ok := msg["id"].(string); ok {
            // Send to broadcaster if from watcher
            if s.broadcaster != nil && s.broadcaster.Conn.RemoteAddr().String() == targetID {
                s.broadcaster.mu.Lock()
                s.broadcaster.WriteJSON(map[string]interface{}{
                    "type": "candidate",
                    "id":   conn.Conn.RemoteAddr().String(),
                    "candidate": msg["candidate"],
                })
                s.broadcaster.mu.Unlock()
            } else {
                // Send to specific watcher if from broadcaster
                for watcher := range s.watchers {
                    if watcher.Conn.RemoteAddr().String() == targetID {
                        watcher.mu.Lock()
                        watcher.WriteJSON(map[string]interface{}{
                            "type": "candidate",
                            "id":   conn.Conn.RemoteAddr().String(),
                            "candidate": msg["candidate"],
                        })
                        watcher.mu.Unlock()
                        break
                    }
                }
            }
        }
    }
}

func (s *Server) handleDisconnect(conn *Connection) {
    s.mu.Lock()
    defer s.mu.Unlock()

    if conn == s.broadcaster {
        s.broadcaster = nil
        // Notify all watchers that broadcaster disconnected
        for watcher := range s.watchers {
            watcher.mu.Lock()
            watcher.WriteJSON(map[string]string{"type": "broadcasterDisconnected"})
            watcher.mu.Unlock()
        }
    } else {
        delete(s.watchers, conn)
        delete(s.waitingWatchers, conn)
        if s.broadcaster != nil {
            s.broadcaster.mu.Lock()
            s.broadcaster.WriteJSON(map[string]interface{}{
                "type": "disconnectPeer",
                "id":   conn.Conn.RemoteAddr().String(),
            })
            s.broadcaster.mu.Unlock()
        }
    }
    conn.Close()
}

func main() {
    server := NewServer()

    // Create HTTP mux for each port
    mux3000 := http.NewServeMux()
    mux3001 := http.NewServeMux()

    // Set up WebSocket handler for both muxes
    mux3000.HandleFunc("/ws", server.handleWebSocket)
    mux3001.HandleFunc("/ws", server.handleWebSocket)

    // Set up different root handlers for each port
    mux3000.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/" {
            http.ServeFile(w, r, "public/listen.html")
            return
        }
        // Serve other static files from public directory
        http.FileServer(http.Dir("public")).ServeHTTP(w, r)
    })

    mux3001.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/" {
            http.ServeFile(w, r, "public/broadcast.html")
            return
        }
        // Serve other static files from public directory
        http.FileServer(http.Dir("public")).ServeHTTP(w, r)
    })

    // Start HTTP server
    go func() {
        log.Printf("Starting HTTP server on :3000")
        if err := http.ListenAndServe(":3000", mux3000); err != nil {
            log.Fatal("HTTP Server error:", err)
        }
    }()

    // Start HTTPS server
    log.Printf("Starting HTTPS server on :3001")
    if err := http.ListenAndServeTLS(":3001", "cert.pem", "key.pem", mux3001); err != nil {
        log.Fatal("HTTPS Server error:", err)
    }
}
