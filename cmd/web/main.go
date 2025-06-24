package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	 _ "github.com/go-sql-driver/mysql"
	 "github.com/joho/godotenv"
)

type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
}

func main() {
	// Importing dotenv
	err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
	port := flag.String("port",":8000","HTTP Network Address") 
	dsn := flag.String("dsn",os.Getenv("MYSQL_DSN"),"MySQL data source name")

	flag.Parse() // Parsing args

	/* Custom Loggers
	Writing to stdout allows to write logs to a file ( go run ./cmd/web >>/tmp/info.log 2>>/tmp/error.lo) */
	infoLog := log.New(os.Stdout,"INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout,"ERROR\t", log.Ldate|log.Ltime)

	// Database connection
	db,err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

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
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB,error) {
	db,err := sql.Open("mysql",dsn)
	if err != nil {
		return nil,err
	}
	if err = db.Ping();err != nil {
		return nil,err
	}
	return db,nil
} 