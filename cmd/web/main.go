package main

import (
	"database/sql"
	"example.com/practice-rest/internal/models"
	"flag"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
	"time"
)

// TODO - create a struct to hold application-wide dependencies
type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	debugLog      *log.Logger
	httpLog       *log.Logger
	post          *models.PostModel
	sessionManger *scs.SessionManager
}

// With http.NewServeMux()
func main() {
	// flag is to define a command line flag, so then we are passing like this -addr=":5000"
	addr := flag.String("addr", "localhost:5000", "HTTP network address to start the server")
	dsn := flag.String("dsn", "root:root@/go_practice?parseTime=true", "MySQL data source name")

	// Parse parses the command-line flags from os.Args[1:]. Must be called after all flags are defined and before flags are accessed by the program.
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	debugLog := log.New(os.Stdout, "DEBUG\t", log.Ldate|log.Ltime|log.Lshortfile)
	httpLog := log.New(os.Stdout, "HTTP\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	// Initializing the session manager using cookies for now,
	// later I'll use jwt to manage the session
	sessionManger := scs.New()
	sessionManger.Store = mysqlstore.New(db)
	sessionManger.Lifetime = 12 * time.Hour

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		debugLog:      debugLog,
		httpLog:       httpLog,
		post:          &models.PostModel{DB: db},
		sessionManger: sessionManger,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Server started on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
