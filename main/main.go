package main

import (
	"database/sql"
	"fmt"
	"http-practive/internal/handler"
	"http-practive/internal/storge"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	db,err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil{
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil{
		log.Fatal(err)
	}
	store := storge.NewPostgresTaskStore(db)
	h := handler.NewHandler(store)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks",h.Get)
	mux.HandleFunc("GET /tasks/{id}",h.TaskByID)
	mux.HandleFunc("POST /tasks",h.TaskCreate)
	mux.HandleFunc("PATCH /tasks/{id}",h.Update)
	mux.HandleFunc("/ping",handler.HandlePing)
	mux.HandleFunc("/hello",handler.HandleHello)
	mux.HandleFunc("/", handler.HandleNotFound)
	
	fmt.Println("Севрвер запущен на порту:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil{
		log.Fatal(err)
	}

}