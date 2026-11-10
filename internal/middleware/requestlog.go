package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"api-template/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const responseMessageCap = 2048

type cappedResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
	max  int
}

func (w *cappedResponseWriter) Write(b []byte) (int, error) {
	if w.body.Len() < w.max {
		remain := w.max - w.body.Len()
		if len(b) > remain {
			w.body.Write(b[:remain])
		} else {
			w.body.Write(b)
		}
	}
	return w.ResponseWriter.Write(b)
}

func RequestLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-Id")
		if requestID == "" {
			requestID = uuid.NewString()
		}
		c.Set("requestId", requestID)
		c.Header("X-Request-Id", requestID)
		c.Request = c.Request.WithContext(
			context.WithValue(c.Request.Context(), "requestId", requestID),
		)

		rw := &cappedResponseWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
			max:            responseMessageCap,
		}
		c.Writer = rw

		start := time.Now()
		c.Next()

		status := c.Writer.Status()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		level := slog.LevelInfo
		switch {
		case status >= http.StatusInternalServerError:
			level = slog.LevelError
		case status >= http.StatusBadRequest:
			level = slog.LevelWarn
		}

		msg := apiResponseMessage(rw.body.Bytes())
		if msg == "" {
			msg = http.StatusText(status)
			if msg == "" {
				msg = "request"
			}
		}

		attrs := []any{
			slog.String("requestId", requestID),
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.Int("status", status),
			slog.Int64("latency_ms", time.Since(start).Milliseconds()),
		}
		if employeeID, ok := c.Get("employeeId"); ok {
			if id, ok := employeeID.(string); ok && id != "" {
				attrs = append(attrs, slog.String("employeeId", id))
			}
		}

		logger.L().Log(c.Request.Context(), level, msg, attrs...)
	}
}

func apiResponseMessage(body []byte) string {
	body = bytes.TrimSpace(body)
	if len(body) == 0 || body[0] != '{' {
		return ""
	}
	dec := json.NewDecoder(bytes.NewReader(body))
	tok, err := dec.Token()
	if err != nil || tok != json.Delim('{') {
		return ""
	}
	for dec.More() {
		key, err := dec.Token()
		if err != nil {
			return ""
		}
		k, _ := key.(string)
		if k == "message" {
			var msg string
			if err := dec.Decode(&msg); err != nil {
				return ""
			}
			return msg
		}
		var skip any
		if err := dec.Decode(&skip); err != nil {
			return ""
		}
	}
	return ""
}

func Recover() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		logger.L().Error("panic",
			"requestId", requestID(c),
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"err", fmt.Sprint(recovered),
		)
		c.AbortWithStatus(http.StatusInternalServerError)
	})
}

func requestID(c *gin.Context) string {
	v, _ := c.Get("requestId")
	id, _ := v.(string)
	return id
}

