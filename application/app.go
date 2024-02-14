package application

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	router http.Handler
	db     *sql.DB
}

func New() *App {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		panic(err)
	}
	app := &App{
		db: db,
	}
	app.loadRoutes()
	return app
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":3000",
		Handler: a.router,
	}

	err := a.db.Ping()

	statement, err := a.db.Prepare(`CREATE TABLE IF NOT EXISTS` + "`order`" + `(
        order_id INTEGER PRIMARY KEY, 
		image TEXT,
        customer_id INTEGER, 
        line_items TEXT, 
        created_at DATETIME,
        shipped_at DATETIME,
        completed_at DATETIME)`)
	if err != nil {
		panic(err)
	}
	statement.Exec()

	// statement, err = a.db.Prepare(`CREATE TABLE IF NOT EXISTS line_item (
	//     item_id INTEGER PRIMARY KEY,
	//     quantity INTEGER,
	//     price INTEGER)`)
	// if err != nil {
	// 	panic(err)
	// }
	// statement.Exec()

	defer func() {
		if err := a.db.Close(); err != nil {
			fmt.Println("failes to close ", err)
		}
	}()

	fmt.Println("Starting server")

	ch := make(chan error, 1)

	go func() {
		err = server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server: %w", err)
		}
		close(ch)
	}()

	select {
	case err = <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		return server.Shutdown(timeout)
	}

}
