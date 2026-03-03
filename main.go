package main

import (
	"database/sql"
	"log"
	"net"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "file:series.db")
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	initDB(db)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
	defer listener.Close()

	log.Println("Servidor corriendo en http://localhost:8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go handleClient(conn, db)
	}
}
