package api

import (
	"context"
	"digital.vasic.translator/pkg/logger"
	"digital.vasic.translator/pkg/translator"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Server represents the API server
type Server struct {
	config     ServerConfig
	router     *gin.Engine
	translator translator.Translator
}

// ServerConfig holds the server configuration
type ServerConfig struct {
	Port     int
	Logger   logger.Logger
	Security *SecurityConfig
}

// SecurityConfig holds security-related configuration
type SecurityConfig struct {
	APIKey         string
	RequireAuth    bool
	MaxRequestSize int64
	MaxBatchSize   int
	RateLimit      int
	RateWindow     time.Duration
	EnableCSRF     bool
	SanitizeInput  bool
	MaxTextLength  int
}

// NewServer creates a new API server
func NewServer(config ServerConfig) *Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	
	server := &Server{
		config: config,
		router: router,
	}
	
	// Add middleware
	router.Use(gin.Recovery())
	if config.Logger != nil {
		router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
			config.Logger.Info("API request", map[string]interface{}{
				"path":     param.Request.URL.Path,
				"method":   param.Request.Method,
				"status":   param.StatusCode,
				"latency":  param.Latency,
				"client_ip": param.ClientIP,
			})
			return ""
		}))
	}
	
	// Add security middleware
	if config.Security != nil {
		router.Use(server.authMiddleware())
	}
	
	// Register routes
	server.registerRoutes()
	
	return server
}

// GetRouter returns the gin router
func (s *Server) GetRouter() *gin.Engine {
	return s.router
}

// Start starts the server
func (s *Server) Start(ctx context.Context) error {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.Port),
		Handler: s.router,
	}
	
	return srv.ListenAndServe()
}

// Stop stops the server
func (s *Server) Stop(ctx context.Context) error {
	// Implementation would need to track the server instance
	// For now, this is a placeholder
	return nil
}

// SetTranslator sets the translator implementation
func (s *Server) SetTranslator(t translator.Translator) {
	s.translator = t
}

// authMiddleware handles authentication
func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if s.config.Security == nil || !s.config.Security.RequireAuth {
			c.Next()
			return
		}
		
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API key required"})
			c.Abort()
			return
		}
		
		if apiKey != s.config.Security.APIKey {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// registerRoutes sets up the API routes
func (s *Server) registerRoutes() {
	// Health check
	s.router.GET("/health", s.healthCheck)
	
	// API routes
	api := s.router.Group("/api")
	{
		api.POST("/translate", s.translateHandler)
		api.GET("/languages", s.languagesHandler)
		api.GET("/stats", s.statsHandler)
		api.POST("/upload", s.uploadHandler)
		api.POST("/batch", s.batchHandler)
	}
}

// Handler functions
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"translator": func() string {
			if s.translator != nil {
				return s.translator.GetName()
			}
			return "none"
		}(),
	})
}

func (s *Server) translateHandler(c *gin.Context) {
	var req struct {
		Text        string `json:"text" binding:"required"`
		SourceLang  string `json:"source_lang" binding:"required"`
		TargetLang  string `json:"target_lang" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if s.translator == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Translator not available"})
		return
	}
	
	// Create context string for the translator
	contextStr := req.SourceLang + "->" + req.TargetLang
	result, err := s.translator.Translate(c.Request.Context(), req.Text, contextStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"translated_text": result,
		"source_lang":     req.SourceLang,
		"target_lang":     req.TargetLang,
	})
}

func (s *Server) languagesHandler(c *gin.Context) {
	// Return supported languages
	c.JSON(http.StatusOK, gin.H{
		"languages": []string{
			"en", "es", "fr", "de", "it", "pt", "ru", "ja", "ko", "zh",
		},
	})
}

func (s *Server) statsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"translations": 0,
		"uptime":       "0s",
	})
}

func (s *Server) uploadHandler(c *gin.Context) {
	// Handle file upload
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

func (s *Server) batchHandler(c *gin.Context) {
	// Handle batch translation
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}