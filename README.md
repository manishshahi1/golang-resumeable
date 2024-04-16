# Chat Application Setup

## Prerequisites
- Kafka installed on your system
- Frontend application running on `localhost:8080`

## Setting up Kafka Topic
1. Open Command Prompt or Terminal.
2. Navigate to the Kafka directory.
3. Run the following command to create a Kafka topic named `chatz`:
    ```bash
    bin\windows\kafka-topics.bat --create --topic chatz --partitions 3 --replication-factor 1 --bootstrap-server localhost:9092 
    ```
   This command will create a topic named `chatz` with 3 partitions and a replication factor of 1.

## Frontend Setup
If the frontend generates the wrong client ID, follow these steps:

1. Open a web browser.
2. Enter the following URL in the address bar:
```http://localhost:8080/fetchPendingMessages?clientId=2```

Replace `2` with the correct client ID of the disconnected websocket.

## Troubleshooting
- If you encounter any issues during setup, refer to the Kafka documentation or the frontend application's documentation for troubleshooting steps.
- Ensure that all dependencies are installed and configured correctly before running the application.

## Additional Notes
- Make sure Kafka is running and accessible on `localhost:9092`.
- Verify that the frontend application is configured to connect to the correct Kafka server and topic.
