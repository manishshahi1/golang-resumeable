// [ DO THIS FIRST ]
// bin\windows\kafka-topics.bat --create --topic chatz --partitions 3 --replication-factor 1 --bootstrap-server localhost:9092 
// Create topic in windows
// Frontend generates wrong Client ID: use this - http://localhost:8080/fetchPendingMessages?clientId=2 (where 2 is the client ID who was disconnected websocket)
