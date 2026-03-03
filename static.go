package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func handleStatic(conn net.Conn, path string) {
	// Build local file path from URL path
	filePath := "." + path // e.g. ./static/style.css

	data, err := os.ReadFile(filePath)
	if err != nil {
		send404(conn)
		return
	}

	// Set content type based on file extension
	contentType := "text/plain"
	if strings.HasSuffix(path, ".css") {
		contentType = "text/css; charset=utf-8"
	} else if strings.HasSuffix(path, ".js") {
		contentType = "application/javascript; charset=utf-8"
	} else if strings.HasSuffix(path, ".svg") {
		contentType = "image/svg+xml"
	}

	response := fmt.Sprintf(
		"HTTP/1.1 200 OK\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n",
		contentType, len(data),
	)
	conn.Write([]byte(response))
	conn.Write(data)
}
