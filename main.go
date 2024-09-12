package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/fatih/color"
)

var (
	clients      = make(map[net.Conn]string)
	clientColors = make(map[net.Conn]*color.Color)
	mu           sync.Mutex
	colors       = []*color.Color{
		color.New(color.FgRed),
		color.New(color.FgGreen),
		color.New(color.FgYellow),
		color.New(color.FgBlue),
		color.New(color.FgMagenta),
		color.New(color.FgCyan),
	}
)

func printBanner() {
	banner := `
'########::'######::'##::::'##::'#######:::'######::'########::'##::::'##:'########:'########::'########:
 ##.....::'##... ##: ##:::: ##:'##.... ##:'##... ##: ##.... ##: ##:::: ##: ##.....:: ##.... ##: ##.....::
 ##::::::: ##:::..:: ##:::: ##: ##:::: ##: ##:::..:: ##:::: ##: ##:::: ##: ##::::::: ##:::: ##: ##:::::::
 ######::: ##::::::: #########: ##:::: ##:. ######:: ########:: #########: ######::: ########:: ######:::
 ##...:::: ##::::::: ##.... ##: ##:::: ##::..... ##: ##.....::: ##.... ##: ##...:::: ##.. ##::: ##...::::
 ##::::::: ##::: ##: ##:::: ##: ##:::: ##:'##::: ##: ##:::::::: ##:::: ##: ##::::::: ##::. ##:: ##:::::::
 ########:. ######:: ##:::: ##:. #######::. ######:: ##:::::::: ##:::: ##: ########: ##:::. ##: ########:
........:::......:::..:::::..:::.......::::......:::..:::::::::..:::::..::........::..:::::..::........::
`
	fmt.Println(color.GreenString(banner))
}

func main() {
	printBanner() // Print the banner when the server starts

	// Set up signal handling for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-stop
		log.Println("Shutting down server...")
		shutdownServer()
		os.Exit(0)
	}()

	port := "2000" // Default port

	// Start server
	log.Printf("Starting server on port %s...", port)
	listenAndServe(port)
}

func listenAndServe(port string) {
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	log.Printf("Server started on port %s", port)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		log.Println("New client connected:", conn.RemoteAddr())

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	// Assign a random color to the client
	rand.Seed(time.Now().UnixNano())
	clientColor := colors[rand.Intn(len(colors))]
	clientColors[conn] = clientColor

	// Prompt user for username
	fmt.Fprint(conn, color.CyanString("\nWelcome! Please enter your username: "))
	username, _ := bufio.NewReader(conn).ReadString('\n')
	username = strings.TrimSpace(username)

	// Add the new client to the map with their username
	mu.Lock()
	clients[conn] = username
	mu.Unlock()

	// Send welcome message and instructions
	welcomeMsg := fmt.Sprintf("\nWelcome to the chat, %s!\n", username)
	instructions := "Type /exit to leave the chat room.\n"
	fmt.Fprint(conn, color.GreenString(welcomeMsg))
	fmt.Fprint(conn, color.CyanString(instructions))
	updateClientList()

	buf := bufio.NewReader(conn)
	for {
		message, err := buf.ReadString('\n')
		if err != nil {
			log.Println("Client disconnected:", conn.RemoteAddr())
			removeClient(conn)
			return
		}

		message = strings.TrimSpace(message)
		if message == "/exit" {
			fmt.Fprint(conn, color.RedString("\nYou have left the chat room.\n"))
			removeClient(conn)
			return
		}

		formattedMessage := fmt.Sprintf("%s: %s", username, message)
		colorMessage := clientColor.Sprintf(formattedMessage + "\n")
		broadcastMessage([]byte(colorMessage), conn)

		// Show the message to the sender as well
		fmt.Fprint(conn, colorMessage)
	}
}

func broadcastMessage(message []byte, sender net.Conn) {
	mu.Lock()
	defer mu.Unlock()

	for client := range clients {
		if client != sender {
			_, err := client.Write(message)
			if err != nil {
				log.Println("Error sending message to:", client.RemoteAddr())
			}
		}
	}
}

func updateClientList() {
	mu.Lock()
	defer mu.Unlock()

	clientList := "\nCurrently in the chat room:\n"
	for conn := range clients {
		clientList += fmt.Sprintf("  - %s\n", clients[conn])
	}

	for conn := range clients {
		fmt.Fprint(conn, color.CyanString(clientList))
	}
}

func removeClient(conn net.Conn) {
	mu.Lock()
	defer mu.Unlock()

	username := clients[conn]
	delete(clients, conn)
	delete(clientColors, conn)

	// Notify other clients about the user leaving
	leaveMessage := fmt.Sprintf("%s has left the chat room.\n", username)
	colorMessage := color.RedString(leaveMessage)
	broadcastMessage([]byte(colorMessage), nil)

	updateClientList()
}

func shutdownServer() {
	mu.Lock()
	defer mu.Unlock()

	for conn := range clients {
		fmt.Fprint(conn, color.RedString("\nServer is shutting down. Goodbye!\n"))
		conn.Close()
	}
}
