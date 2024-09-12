# Echo$phere Chat Server

Echo$phere is a simple chat server implemented in Go, providing real-time text-based communication between multiple clients. It features color-coded messages for a better user experience and supports graceful shutdowns and client management.

## Features

- **Multi-client Support**: Handle multiple clients concurrently.
- **Color-coded Messages**: Messages are color-coded for better readability.
- **Username Management**: Users provide a username when joining.
- **Client List**: Clients see a list of currently active users.
- **Graceful Shutdown**: The server handles interruptions and shuts down gracefully.
- **Exit Command**: Clients can exit the chat room using a specific command.

## Requirements

- **Go**: Version 1.18 or higher.

## Installation

### 1. **Clone the Repository**

   ```bash
   git clone https://github.com/yourusername/echosphere-chat-server.git
   cd echo$phere-chat-server
   ```

### 2. Install Dependencies

The project uses the fatih/color package for color handling. Install it using:

     ```bash
     go get github.com/fatih/color
     ```

### 3. Build the Project

Compile the Go code into an executable:

     ```bash
     go build -o server main.go
     ```

### 4. Usage
  Run the server executable to start the server:

     ```bash
     ./server
     ```

  The server will start and listen on port 2000 by default. You can change the port by modifying the port variable in main.go.

### 5. Connect to the Server

  Use a TCP client (like telnet or a custom client) to connect to the server:

     ```bash
     telnet localhost 2000
     ```

  Or, if you have a custom client implementation:
     ```bash
     go run client.go
     ```

### 7. Interact with the Chat Room

  Enter a Username: You'll be prompted to enter a username upon connection.
  Send Messages: Type your message and press Enter to send it.
  Exit the Chat Room: Type /exit to leave the chat room.
  Code Structure
  The project consists of a single file main.go, which handles both server and client operations:
  
  main.go: Contains the chat server implementation, including client management, message broadcasting, and graceful shutdowns.
  Contributing
  Feel free to open issues or submit pull requests to contribute to the project. Please ensure that your code adheres to the 
  project's coding standards and pass all tests.
  
