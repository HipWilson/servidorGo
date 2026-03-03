package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/url"
	"strconv"
	"strings"
)

func handleClient(conn net.Conn, db *sql.DB) {
	defer conn.Close()

	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil || n == 0 {
		return
	}
	request := string(buffer[:n])

	// Parse first line: METHOD /path HTTP/1.1
	lines := strings.Split(request, "\r\n")
	parts := strings.Fields(lines[0])
	if len(parts) < 2 {
		return
	}
	method := parts[0]
	rawPath := parts[1]

	// Separate path from query string
	pathParts := strings.SplitN(rawPath, "?", 2)
	path := pathParts[0]
	queryString := ""
	if len(pathParts) > 1 {
		queryString = pathParts[1]
	}

	// Read Content-Length header
	contentLength := 0
	for _, line := range lines[1:] {
		if strings.HasPrefix(strings.ToLower(line), "content-length:") {
			lengthStr := strings.TrimSpace(line[len("content-length:"):])
			contentLength, _ = strconv.Atoi(lengthStr)
		}
	}

	// Read body (comes after blank line \r\n\r\n)
	body := ""
	if idx := strings.Index(request, "\r\n\r\n"); idx != -1 {
		body = request[idx+4:]
		if contentLength > 0 && len(body) > contentLength {
			body = body[:contentLength]
		}
	}

	log.Printf("%s %s", method, rawPath)

	// Routes
	if method == "GET" && path == "/" {
		handleIndex(conn, db)
	} else if method == "GET" && path == "/create" {
		handleCreateForm(conn)
	} else if method == "POST" && path == "/create" {
		handleCreatePost(conn, db, body)
	} else if method == "POST" && path == "/update" {
		params, _ := url.ParseQuery(queryString)
		id := params.Get("id")
		handleUpdate(conn, db, id)
	} else if method == "POST" && path == "/decrement" {
		params, _ := url.ParseQuery(queryString)
		id := params.Get("id")
		handleDecrement(conn, db, id)
	} else if method == "DELETE" && path == "/delete" {
		params, _ := url.ParseQuery(queryString)
		id := params.Get("id")
		handleDelete(conn, db, id)
	} else if method == "GET" && strings.HasPrefix(path, "/static/") {
		handleStatic(conn, path)
	} else {
		send404(conn)
	}
}

func handleIndex(conn net.Conn, db *sql.DB) {
	seriesList, err := getAllSeries(db)
	if err != nil {
		log.Println("Error getting series:", err)
		send500(conn)
		return
	}

	rows := ""
	for _, s := range seriesList {
		// Calculate progress percentage
		pct := 0
		if s.TotalEpisodes > 0 {
			pct = (s.CurrentEpisode * 100) / s.TotalEpisodes
		}

		// Mark completed series
		status := ""
		if s.CurrentEpisode >= s.TotalEpisodes {
			status = " (Completada)"
		}

		// Disable buttons when at limits
		decDisabled := ""
		if s.CurrentEpisode <= 1 {
			decDisabled = "disabled"
		}
		incDisabled := ""
		if s.CurrentEpisode >= s.TotalEpisodes {
			incDisabled = "disabled"
		}

		rows += fmt.Sprintf(
			"<tr data-id=\"%d\">"+
				"<td>%d</td>"+
				"<td>%s%s</td>"+
				"<td>%d / %d</td>"+
				"<td><div class=\"progress-bar-bg\"><div class=\"progress-bar-fill\" style=\"width:%d%%\"></div></div> %d%%</td>"+
				"<td>"+
				"<button onclick=\"changeEp(%d,'decrement')\" %s>-1</button>"+
				"<button onclick=\"changeEp(%d,'update')\" %s>+1</button>"+
				"<button onclick=\"deleteSeries(%d)\">Eliminar</button>"+
				"</td>"+
				"</tr>",
			s.ID,
			s.ID,
			s.Name, status,
			s.CurrentEpisode, s.TotalEpisodes,
			pct, pct,
			s.ID, decDisabled,
			s.ID, incDisabled,
			s.ID,
		)
	}

	html := buildIndexPage(rows)
	sendHTML(conn, 200, html)
}

func handleCreateForm(conn net.Conn) {
	sendHTML(conn, 200, buildCreatePage("", "", "1", ""))
}

func handleCreatePost(conn net.Conn, db *sql.DB, body string) {
	values, err := url.ParseQuery(body)
	if err != nil {
		sendText(conn, 400, "Bad Request")
		return
	}

	name := strings.TrimSpace(values.Get("series_name"))
	currentStr := strings.TrimSpace(values.Get("current_episode"))
	totalStr := strings.TrimSpace(values.Get("total_episodes"))

	// Server-side validation
	errorMsg := ""
	if name == "" {
		errorMsg = "El nombre no puede estar vacio."
	}

	current, err := strconv.Atoi(currentStr)
	if err != nil || current < 1 {
		errorMsg = "El episodio actual debe ser un numero valido mayor a 0."
	}

	total, err := strconv.Atoi(totalStr)
	if err != nil || total < 1 {
		errorMsg = "El total de episodios debe ser un numero valido mayor a 0."
	}

	if errorMsg == "" && current > total {
		errorMsg = "El episodio actual no puede ser mayor que el total de episodios."
	}

	if errorMsg != "" {
		sendHTML(conn, 400, buildCreatePage(errorMsg, name, currentStr, totalStr))
		return
	}

	err = insertSeries(db, name, current, total)
	if err != nil {
		log.Println("Error inserting series:", err)
		send500(conn)
		return
	}

	// POST/Redirect/GET pattern
	conn.Write([]byte("HTTP/1.1 303 See Other\r\nLocation: /\r\n\r\n"))
}

func handleUpdate(conn net.Conn, db *sql.DB, id string) {
	if id == "" {
		sendText(conn, 400, "missing id")
		return
	}
	err := incrementEpisode(db, id)
	if err != nil {
		log.Println("Error updating episode:", err)
		sendText(conn, 500, "error")
		return
	}
	sendText(conn, 200, "ok")
}

func handleDecrement(conn net.Conn, db *sql.DB, id string) {
	if id == "" {
		sendText(conn, 400, "missing id")
		return
	}
	err := decrementEpisode(db, id)
	if err != nil {
		log.Println("Error decrementing episode:", err)
		sendText(conn, 500, "error")
		return
	}
	sendText(conn, 200, "ok")
}

func handleDelete(conn net.Conn, db *sql.DB, id string) {
	if id == "" {
		sendText(conn, 400, "missing id")
		return
	}
	err := deleteSeries(db, id)
	if err != nil {
		log.Println("Error deleting series:", err)
		sendText(conn, 500, "error")
		return
	}
	sendText(conn, 200, "ok")
}

func sendHTML(conn net.Conn, status int, html string) {
	statusText := "OK"
	if status == 400 {
		statusText = "Bad Request"
	}
	response := fmt.Sprintf(
		"HTTP/1.1 %d %s\r\nContent-Type: text/html; charset=utf-8\r\nContent-Length: %d\r\n\r\n%s",
		status, statusText, len(html), html,
	)
	conn.Write([]byte(response))
}

func sendText(conn net.Conn, status int, text string) {
	statusText := "OK"
	if status != 200 {
		statusText = "Error"
	}
	response := fmt.Sprintf(
		"HTTP/1.1 %d %s\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s",
		status, statusText, len(text), text,
	)
	conn.Write([]byte(response))
}

func send404(conn net.Conn) {
	body := "<h1>404 - Pagina no encontrada</h1>"
	response := fmt.Sprintf(
		"HTTP/1.1 404 Not Found\r\nContent-Type: text/html\r\nContent-Length: %d\r\n\r\n%s",
		len(body), body,
	)
	conn.Write([]byte(response))
}

func send500(conn net.Conn) {
	body := "500 - Error interno del servidor"
	response := fmt.Sprintf(
		"HTTP/1.1 500 Internal Server Error\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s",
		len(body), body,
	)
	conn.Write([]byte(response))
}
