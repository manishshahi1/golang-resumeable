package main

import (
    "context"
    "log"
    "net/http"
    "strconv"
    "sync"

    "github.com/gorilla/websocket"
    "github.com/segmentio/kafka-go"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

var (
    clients      = make(map[*websocket.Conn]string)
    clientsMu    sync.Mutex
    clientCounter int
)

func getNextClientId() string {
    clientsMu.Lock()
    defer clientsMu.Unlock()
    clientCounter++
    return strconv.Itoa(clientCounter)
}

func handleMessages(w http.ResponseWriter, r *http.Request) {
    // Upgrade HTTP connection to WebSocket
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }
    defer conn.Close()

    // Get client ID
    clientId := getNextClientId()
    log.Printf("Client %s connected\n", clientId)

    clientsMu.Lock()
    clients[conn] = clientId
    clientsMu.Unlock()
    defer func() {
        clientsMu.Lock()
        delete(clients, conn)
        clientsMu.Unlock()
        log.Printf("Client %s disconnected\n", clientId)
    }()

    // Kafka setup
    writer := kafka.NewWriter(kafka.WriterConfig{
        Brokers: []string{"localhost:9092"},
        Topic:   "chatz", // Set the topic name as "chatz"
    })
    defer writer.Close()

    for {
        // Read message from WebSocket client
        _, msg, err := conn.ReadMessage()
        if err != nil {
            log.Printf("Client %s: %v\n", clientId, err)
            break
        }

        // Publish message to Kafka topic
        err = writer.WriteMessages(context.Background(), kafka.Message{
            Key:   []byte(clientId), // Use client ID as the key
            Value: msg,
        })
        if err != nil {
            log.Printf("Client %s: Failed to publish message to Kafka: %v\n", clientId, err)
        }

        // Broadcast message to all connected WebSocket clients
        clientsMu.Lock()
        for client := range clients {
            err = client.WriteMessage(websocket.TextMessage, msg)
            if err != nil {
                log.Printf("Client %s: Failed to send message to client: %v\n", clientId, err)
                return
            }
        }
        clientsMu.Unlock()
    }
}

func fetchPendingMessages(w http.ResponseWriter, r *http.Request) {
    // Get client ID from query parameters
    clientId := r.URL.Query().Get("clientId")

    // Kafka setup
    reader := kafka.NewReader(kafka.ReaderConfig{
        Brokers:  []string{"localhost:9092"},
        Topic:    "chatz", // Set the topic name as "chatz"
        GroupID:  "my-consumer-group",
        MinBytes: 10e3, // 10KB
        MaxBytes: 10e6, // 10MB
    })
    defer reader.Close()

    // Fetch pending messages
    for {
        m, err := reader.ReadMessage(context.Background())
        if err != nil {
            log.Println("Error reading message:", err)
            break
        }

        // Send the message back to the client
        clientsMu.Lock()
        for client, id := range clients {
            if id == clientId {
                err := client.WriteMessage(websocket.TextMessage, m.Value)
                if err != nil {
                    log.Printf("Client %s: Failed to send message to client: %v\n", clientId, err)
                    // Handle error, maybe reconnect client or mark message as pending
                }
                break
            }
        }
        clientsMu.Unlock()
    }
}

func main() {
    // Define WebSocket and message fetching handlers
    http.HandleFunc("/ws", handleMessages)
    http.HandleFunc("/fetchPendingMessages", fetchPendingMessages)

    // Serve the msg.html file from the root directory
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "msg.html")
    })

    // Start the HTTP server
    log.Fatal(http.ListenAndServe(":8080", nil))
}


// bin\windows\kafka-topics.bat --create --topic chatz --partitions 3 --replication-factor 1 --bootstrap-server localhost:9092 
// Create topic in windows
// Frontend generates wrong Client ID: use this - http://localhost:8080/fetchPendingMessages?clientId=2 (where 2 is the client ID who was disconnected websocket)
