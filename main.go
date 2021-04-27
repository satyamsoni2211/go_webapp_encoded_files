package main

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
	"github.com/satyamsoni2211/go_webapp_encoded_files/internal/handler"
)

var (
	templates *template.Template
	err       error
)

func init() {
	// templates, err = template.ParseFiles("templates/*")
}
func getTemplate(name string) (*template.Template, error) {
	templates, err = template.ParseFiles(fmt.Sprintf("templates/%s", name))
	if err != nil {
		return nil, err
	}
	return templates, nil
}

func rootRoute(w http.ResponseWriter, r *http.Request) {
	tpl, err := handler.HandleTemplate("index.html")
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}
	tpl.Execute(w, struct{ Name string }{Name: "Golang"})
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", rootRoute).Methods("GET")
	fileHandler := handler.NewFileHandler()
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileHandler))
	defer func() {
		fmt.Println("Cleaning all temp files from ", handler.Dir)
		os.RemoveAll(handler.Dir)
	}()
	server := &http.Server{
		Addr:    "[::]:5000",
		Handler: router,
	}

	go func() {
		err := server.ListenAndServe()
		//Checking if error is not due to server closed
		if err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				fmt.Printf("Error Occurred while startign server %s", err)
			}
		}
	}()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Println("Received interrupt signal")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := server.Shutdown(ctx)
	if err != nil {
		fmt.Printf("Error occurred while shutting down")
	}
}
