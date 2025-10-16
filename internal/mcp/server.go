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

	// Add logging middleware (must be added before registering tools)
	mcpServer.AddReceivingMiddleware(LoggingMiddleware(logger))

	// Create and register search tool
	searchTool := tools.NewSearchTool(tmdbClient, logger)
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        searchTool.Name(),
		Description: searchTool.Description(),
	}, searchTool.Handler())

	// Create and register get_details tool
	getDetailsTool := tools.NewGetDetailsTool(tmdbClient, logger)
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        getDetailsTool.Name(),
		Description: getDetailsTool.Description(),
	}, getDetailsTool.Handler())

	// Create and register discover_movies tool
	discoverMoviesTool := tools.NewDiscoverMoviesTool(tmdbClient, logger)
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        discoverMoviesTool.Name(),
		Description: discoverMoviesTool.Description(),
	}, discoverMoviesTool.Handler())

	// Create and register discover_tv tool
	discoverTVTool := tools.NewDiscoverTVTool(tmdbClient, logger)
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        discoverTVTool.Name(),
		Description: discoverTVTool.Description(),
	}, discoverTVTool.Handler())

	// Create and register get_trending tool
	getTrendingTool := tools.NewGetTrendingTool(tmdbClient, logger)
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        getTrendingTool.Name(),
		Description: getTrendingTool.Description(),
	}, getTrendingTool.Handler())

	// Create and register get_recommendations tool
	getRecommendationsTool := tools.NewGetRecommendationsTool(tmdbClient, logger)
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        getRecommendationsTool.Name(),
		Description: getRecommendationsTool.Description(),
	}, getRecommendationsTool.Handler())

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
