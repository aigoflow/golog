package main

import (
	"crypto/subtle"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (e *Engine) apiKeyMiddleware(expectedKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Support both "Bearer token" and just "token" formats
		token := strings.TrimPrefix(auth, "Bearer ")
		token = strings.TrimSpace(token)

		if token != expectedKey {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (e *Engine) uiPasswordMiddleware(expectedPassword string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for session cookie first
		if cookie, err := c.Cookie("ui-auth"); err == nil && cookie == "authenticated" {
			c.Next()
			return
		}

		// Check for login POST
		if c.Request.Method == "POST" && c.Request.URL.Path == "/ui/login" {
			password := c.PostForm("password")
			if subtle.ConstantTimeCompare([]byte(password), []byte(expectedPassword)) == 1 {
				c.SetCookie("ui-auth", "authenticated", 3600*24, "/", "", false, true)
				c.Redirect(http.StatusFound, "/ui")
				return
			}
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Invalid password"})
			return
		}

		// Show login form
		if c.Request.URL.Path == "/ui" || strings.HasPrefix(c.Request.URL.Path, "/ui/") {
			c.HTML(http.StatusOK, "login.html", gin.H{})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (e *Engine) setupRoutes() *gin.Engine {
	// Set Gin mode based on environment
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}
	
	r := gin.New()
	
	// Add default middleware manually to avoid the warning
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	
	// Add CORS middleware for UI
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})
	
	// Configure trusted proxies
	r.SetTrustedProxies([]string{"127.0.0.1"})

	// Load HTML templates
	e.loadTemplates(r)

	// API routes
	api := r.Group("/api/v1")
	
	// Apply API key middleware if API_KEY is set
	apiKey := os.Getenv("API_KEY")
	if apiKey != "" {
		api.Use(e.apiKeyMiddleware(apiKey))
	}
	
	{
		// Session management
		api.POST("/sessions", e.createSessionHandler)
		api.GET("/sessions", e.listSessionsHandler)
		api.GET("/sessions/:id", e.getSessionHandler)
		api.DELETE("/sessions/:id", e.deleteSessionHandler)
		
		// Facts and rules (scoped to sessions)
		api.POST("/sessions/:sessionId/facts", e.addFactHandler)
		api.POST("/sessions/:sessionId/rules", e.addRuleHandler)
		api.POST("/sessions/:sessionId/query", e.queryHandler)
		
		// Cache management
		api.POST("/cache/clear", e.clearCacheHandler)
	}

	// UI routes (if enabled)
	if os.Getenv("ENABLE_UI") == "true" {
		ui := r.Group("/ui")
		
		// Apply UI password middleware if UI_PASSWORD is set
		uiPassword := os.Getenv("UI_PASSWORD")
		if uiPassword != "" {
			ui.Use(e.uiPasswordMiddleware(uiPassword))
		}
		
		// UI login route (must be before middleware)
		r.POST("/ui/login", func(c *gin.Context) {
			// This is handled by the middleware
		})
		
		ui.GET("/", e.uiHandler)
		ui.GET("", e.uiHandler)
		ui.GET("/js", e.jsHandler)
	}

	return r
}

func (e *Engine) createSessionHandler(c *gin.Context) {
	var req CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session, err := e.CreateSession(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, session)
}

func (e *Engine) listSessionsHandler(c *gin.Context) {
	sessions, err := e.ListSessions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"sessions": sessions})
}

func (e *Engine) getSessionHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	session, err := e.GetSession(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	c.JSON(http.StatusOK, session)
}

func (e *Engine) deleteSessionHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	if err := e.DeleteSession(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "session deleted"})
}

func (e *Engine) addFactHandler(c *gin.Context) {
	sessionIdStr := c.Param("sessionId")
	sessionId, err := strconv.Atoi(sessionIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	var fact Fact
	if err := c.ShouldBindJSON(&fact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fact.SessionID = sessionId
	
	// Validate that predicate is not empty
	if fact.Predicate.Type == "" || fact.Predicate.Value == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "predicate is required"})
		return
	}

	if err := e.AddFact(fact); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	e.UpdateSessionTimestamp(sessionId)
	c.JSON(http.StatusOK, gin.H{"status": "fact added"})
}

func (e *Engine) addRuleHandler(c *gin.Context) {
	sessionIdStr := c.Param("sessionId")
	sessionId, err := strconv.Atoi(sessionIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	var rule Rule
	if err := c.ShouldBindJSON(&rule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rule.SessionID = sessionId

	if err := e.AddRule(rule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	e.UpdateSessionTimestamp(sessionId)
	c.JSON(http.StatusOK, gin.H{"status": "rule added"})
}

func (e *Engine) queryHandler(c *gin.Context) {
	sessionIdStr := c.Param("sessionId")
	sessionId, err := strconv.Atoi(sessionIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	var query Query
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := e.Query(query, sessionId)
	e.UpdateSessionTimestamp(sessionId)
	c.JSON(http.StatusOK, result)
}

func (e *Engine) clearCacheHandler(c *gin.Context) {
	e.ClearCache()
	c.JSON(http.StatusOK, gin.H{"status": "cache cleared"})
}

func (e *Engine) loadTemplates(r *gin.Engine) {
	if os.Getenv("ENABLE_UI") == "true" {
		// Load templates from embedded strings
		tmpl := template.New("")
		template.Must(tmpl.New("login.html").Parse(loginTemplate))
		template.Must(tmpl.New("ui.html").Parse(uiTemplate))
		r.SetHTMLTemplate(tmpl)
	}
}

func (e *Engine) uiHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "ui.html", gin.H{
		"Title": "Prolog Engine REPL",
	})
}

func (e *Engine) jsHandler(c *gin.Context) {
	c.Header("Content-Type", "application/javascript")
	c.String(http.StatusOK, jsContent)
}