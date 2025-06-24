package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
}

func main() {
	port := flag.String("port",":8000","HTTP Network Address") 

	flag.Parse() // Parsing args

	/* Custom Loggers
	Writing to stdout allows to write logs to a file ( go run ./cmd/web >>/tmp/info.log 2>>/tmp/error.lo) */
	infoLog := log.New(os.Stdout,"INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout,"ERROR\t", log.Ldate|log.Ltime)

	// Instanciating a new application struct
	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
	}

	/* Initializing a http.Server struct to pass down options */
	srv := &http.Server{
		Addr: *port,
		ErrorLog: errorLog,
		Handler: app.routes(),
	}

	infoLog.Println("Starting server on http://localhost:8000")
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
