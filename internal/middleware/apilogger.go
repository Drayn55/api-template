package middleware

import (
	"bytes"
	"io"
	"math"
	"time"

	"api-template/internal/model"
	"api-template/internal/repository"

	"github.com/gin-gonic/gin"
)

type responseWriter struct {
	gin.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

// APILogger logs every request to the api_logs table if logRepo is provided
func APILogger(logRepo *repository.APILogRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		if logRepo == nil {
			c.Next()
			return
		}

		start := time.Now()

		var requestBody string
		if c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil {
				requestBody = string(bodyBytes)
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		rw := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
			statusCode:     200,
		}
		c.Writer = rw

		c.Next()

		execTime := float64(time.Since(start).Milliseconds())
		execTime = math.Round(execTime*100) / 100

		method := c.Request.Method
		endpoint := c.FullPath()
		responseBody := rw.body.String()
		statusCode := c.Writer.Status()
		if statusCode == 0 {
			statusCode = rw.statusCode
		}
		ipAddress := c.ClientIP()
		userAgent := c.Request.UserAgent()
		employeeID, _ := c.Get("employeeId")

		empID := ""
		if employeeID != nil {
			empID = employeeID.(string)
		}

		go func() {
			_ = logRepo.Create(model.APILog{
				Method:          &method,
				Endpoint:        &endpoint,
				RequestBody:     &requestBody,
				ResponseBody:    &responseBody,
				StatusCode:      &statusCode,
				IPAddress:       &ipAddress,
				UserAgent:       &userAgent,
				EmployeeID:      &empID,
				ExecutionTimeMs: &execTime,
			})
		}()
	}
}
