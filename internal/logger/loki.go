package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type lokiTee struct {
	url         string
	serviceName string
	client      *http.Client
}

func (t *lokiTee) Write(p []byte) (int, error) {
	n, err := os.Stdout.Write(p)
	if err != nil || t.url == "" {
		return n, err
	}
	line := bytes.TrimSpace(p)
	if len(line) == 0 || line[0] != '{' {
		return n, nil
	}
	t.push(line)
	return n, nil
}

func (t *lokiTee) push(line []byte) {
	var rec map[string]any
	if err := json.Unmarshal(line, &rec); err != nil {
		return
	}

	service, _ := rec["service"].(string)
	level, _ := rec["level"].(string)
	if service == "" {
		service = t.serviceName
		if service == "" {
			service = "api-template"
		}
	}
	if level == "" {
		level = "info"
	}

	value := []any{strconv.FormatInt(time.Now().UnixNano(), 10), string(line)}
	if requestID, ok := rec["requestId"].(string); ok && requestID != "" {
		value = append(value, map[string]string{"requestId": requestID})
	}

	body, err := json.Marshal(map[string]any{
		"streams": []any{
			map[string]any{
				"stream": map[string]string{
					"service": service,
					"level":   level,
				},
				"values": []any{value},
			},
		},
	})
	if err != nil {
		return
	}

	if t.client == nil {
		t.client = &http.Client{Timeout: 2 * time.Second}
	}
	req, err := http.NewRequest(http.MethodPost, t.url+"/loki/api/v1/push", bytes.NewReader(body))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "loki push: %v\n", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		fmt.Fprintf(os.Stderr, "loki push: %s %s\n", resp.Status, bytes.TrimSpace(msg))
	}
}
