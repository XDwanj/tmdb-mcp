package mcp

import (
	"context"

	"github.com/XDwanj/tmdb-mcp/internal/tmdb"
	"github.com/XDwanj/tmdb-mcp/internal/tools"
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

	// Register search tool
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "search",
		Description: "Search for movies, TV shows, and people on TMDB using a query string",
	}, func(ctx context.Context, req *mcp.CallToolRequest, params tools.SearchParams) (*mcp.CallToolResult, tools.SearchResponse, error) {
		// Set default page
		if params.Page == 0 {
			params.Page = 1
		}

		logger.Info("Search request received",
			zap.String("query", params.Query),
			zap.Int("page", params.Page),
		)

		// Call TMDB Client (validation is done in the client layer)
		results, err := tmdbClient.Search(ctx, params.Query, params.Page)
		if err != nil {
			logger.Error("Search failed",
				zap.Error(err),
				zap.String("query", params.Query),
			)
			return nil, tools.SearchResponse{}, err
		}

		logger.Info("Search completed",
			zap.String("query", params.Query),
			zap.Int("results", len(results.Results)),
		)

		return &mcp.CallToolResult{}, tools.SearchResponse{Results: results.Results}, nil
	})

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
