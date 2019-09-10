package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"

	"github.com/sergivb01/paste/config"
	"github.com/sergivb01/paste/utils"

	// postgresql driver
	_ "github.com/lib/pq"
)

// Server defines the PasteServer
type Server struct {
	router *mux.Router

	cfg config.Config

	// Postgresql db
	db *sqlx.DB
}

func New(cfgPath string) (*Server, error) {
	c, err := config.LoadFromFile(cfgPath)


	db, err := sqlx.Open("postgres", c.PostgresURI)
	if err != nil {
		return nil, fmt.Errorf("couldn't open postgresql: %v", err)
	}

	if _, err := db.Exec(utils.CreatePastesTable); err != nil {
		return nil, fmt.Errorf("couldn't create pastes table: %w", err)
	}

	s := &Server{
		cfg: c,
		db: db,
	}
	_ = db // ignore declared but not used
	s.routes()

	return s, nil
}

// Listen starts the HTTP server to handle requests
func (s *Server) Listen() {
	srv := &http.Server{
		Addr:         ":8080",
		WriteTimeout: time.Second * 10,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 15,
		Handler:      s.router,
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		fmt.Printf("started listening on %s...\n", ":8080")
		if err := srv.ListenAndServe(); err != nil {
			fmt.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// if err := s.db.Close(); err != nil {
	// 	fmt.Printf("error closing db: %v", err)
	// }

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	srv.Shutdown(ctx)
}
