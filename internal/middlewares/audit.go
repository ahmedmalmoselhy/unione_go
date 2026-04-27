package middlewares

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/ahmedmalmoselhy/unione_go/internal/services"
	"github.com/gin-gonic/gin"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func AuditMiddleware(auditSvc *services.AuditService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only log mutation methods
		method := c.Request.Method
		if method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions {
			c.Next()
			return
		}

		// Skip login and register to avoid logging sensitive data
		path := c.Request.URL.Path
		if strings.Contains(path, "/auth/login") || strings.Contains(path, "/auth/register") {
			c.Next()
			return
		}

		// Read request body
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			// Restore request body for later use
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Wrap writer to capture response body if needed (optional)
		// blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		// c.Writer = blw

		c.Next()

		// Log after request processing
		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			userID := uint(0)
			if id, exists := c.Get("user_id"); exists {
				if uid, ok := id.(uint); ok {
					userID = uid
				}
			}

			// Capture IP and User Agent
			ip := c.ClientIP()
			userAgent := c.Request.UserAgent()

			// Action name from path and method
			action := method + " " + path

			// Entity type and ID can be inferred from path but it's better to be explicit
			// For middleware, we'll just log the path as action
			auditSvc.Log(userID, action, "API_REQUEST", "", string(requestBody), "", ip, userAgent)
		}
	}
}
