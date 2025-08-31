package agent

import (
	"context"
	"log"
	"time"

	"github.com/mark3labs/mcp-go"
)

// MCPAgent handles communication with the MCP server.
type MCPAgent struct {
	client *mcp.Client
	ctx    context.Context
}

// NewMCPAgent creates a new MCP agent with the given server URL.
func NewMCPAgent(serverURL string) (*MCPAgent, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize the MCP client with the server URL and optional options
	client, err := mcp.NewClient(serverURL, mcp.WithTimeout(30*time.Second))
	if err != nil {
		return nil, err
	}

	return &MCPAgent{
		client: client,
		ctx:    ctx,
	}, nil
}

// Start initializes the agent and begins processing.
func (a *MCPAgent) Start() error {
	log.Println("MCP agent starting...")

	// Connect to the server and begin listening for messages
	if err := a.client.Connect(); err != nil {
		return err
	}

	// Subscribe to messages from the server
	go a.listenForMessages()

	log.Println("MCP agent started successfully.")
	return nil
}

// listenForMessages handles incoming messages from the server.
func (a *MCPAgent) listenForMessages() {
	for {
		select {
		case <-a.ctx.Done():
			log.Println("MCP agent shutting down.")
			return
		case msg := <-a.client.Messages():
			if err := a.handleMessage(msg); err != nil {
				log.Printf("Error handling message: %v", err)
			}
		}
	}
}

// handleMessage processes a single message.
func (a *MCPAgent) handleMessage(msg *mcp.Message) error {
	log.Printf("Received message: %v", msg)

	// Example: Echo the message back
	response := &mcp.Message{
		ID:      msg.ID,
		Content: msg.Content,
		Type:    mcp.MessageTypeResponse,
	}

	return a.client.Send(response)
}

// Stop gracefully shuts down the agent.
func (a *MCPAgent) Stop() {
	log.Println("Stopping MCP agent...")
	a.client.Close()
	close(a.ctx.Done())
}

// Ping tests the connection to the server.
func (a *MCPAgent) Ping() error {
	return a.client.Ping()
}