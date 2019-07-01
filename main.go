package main

import (
	"fmt"
	"github.com/benibana2001/hknews_go/Controllers"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		fmt.Println("$PORT must be set, so that set default value")
		port = "3000"
	}

	server := http.Server{
		Addr: ":" + port,
	}

	controller := Controllers.Controller{}
	http.HandleFunc("/", controller.Home)
	server.ListenAndServe()
}
