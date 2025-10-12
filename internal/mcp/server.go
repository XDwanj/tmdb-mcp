package mcp

import (
	"context"

	"github.com/XDwanj/tmdb-mcp/internal/tmdb"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/zap"
)

// Server wraps the MCP server and provides TMDB-specific functionality
type Server struct {
	mcpServer  *mcp.Server
	tmdbClient *tmdb.Client
	logger     *zap.Logger
}

// NewServer creates a new MCP server instance with TMDB client integration
func NewServer(tmdbClient *tmdb.Client, logger *zap.Logger) *Server {
	// Create server options
	opts := &mcp.ServerOptions{
		Instructions: "TMDB Movie Database MCP Server - provides tools for searching and retrieving movie information",
	}

	// Create MCP server with implementation info
	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    "tmdb-mcp",
		Version: "1.0.0",
	}, opts)

	// TODO: Future story 1.6 will add tools here
	// Example: mcp.AddTool(mcpServer, &mcp.Tool{Name: "search"}, searchHandler)

	return &Server{
		mcpServer:  mcpServer,
		tmdbClient: tmdbClient,
		logger:     logger,
	}
}

// Run starts the MCP server with the specified transport
func (s *Server) Run(ctx context.Context, transport mcp.Transport) error {
	s.logger.Info("Starting MCP server")
	return s.mcpServer.Run(ctx, transport)
}
