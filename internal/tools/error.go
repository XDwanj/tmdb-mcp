package tools

import (
	"errors"
	"fmt"

	"github.com/XDwanj/tmdb-mcp/internal/tmdb"
)

// convertTMDBError converts TMDB client errors to user-friendly MCP error messages
func convertTMDBError(err error, resourceType string) error {
	if err == nil {
		return nil
	}

	var tmdbErr *tmdb.TMDBError
	if errors.As(err, &tmdbErr) {
		switch tmdbErr.ErrorType {
		case tmdb.ErrorTypeAuth:
			return fmt.Errorf("authentication failed. Please check your TMDB API Key configuration")
		case tmdb.ErrorTypeNotFound:
			return fmt.Errorf("the requested %s was not found", resourceType)
		case tmdb.ErrorTypeRateLimit:
			return fmt.Errorf("rate limit exceeded. The service will automatically retry shortly")
		case tmdb.ErrorTypeNetwork:
			return fmt.Errorf("request timed out. Please try again or check your network connection")
		case tmdb.ErrorTypeServer:
			return fmt.Errorf("TMDB service is temporarily unavailable. Please try again later")
		case tmdb.ErrorTypeParsing:
			return fmt.Errorf("failed to process TMDB API response. Please try again")
		}
	}

	// 保留原始错误链
	return fmt.Errorf("failed to fetch %s: %w", resourceType, err)
}
