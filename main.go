package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"net/http"
)

type App struct {
	Router *mux.Router
	DB     *pgx.Conn
}

func (a *App) Init() {
	fmt.Println("Initializing...")
	a.Router = mux.NewRouter()
	var err error
	a.DB, err = pgx.Connect(context.Background(), "postgresql://castlemound:d39olrf9e4ov77pc@db-postgresql-lon1-16483-do-user-9489448-0.b.db.ondigitalocean.com:25060/cm?sslmode=require")

	if err != nil {
		fmt.Printf("Error while connecting to the database: %v", err)
	}
	fmt.Println("Successul")
}

func (a *App) Run() {
	fmt.Println("Running the application")
	a.Router.HandleFunc("/categories", a.CategoriesHandler)
	defer a.DB.Close(context.Background())
	fmt.Println("Listening at port 8000")
	http.ListenAndServe(":8000", a.Router)
}

func main() {
	a := App{}
	a.Init()
	a.Run()
}
