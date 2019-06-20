package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type Article struct {
	By          string `json:"by"`
	Descendants int    `json:"descendants"`
	Id          int    `json:"id"`
	Kids        []int  `json:"kids"`
	Score       int    `json:"score"`
	Time        int    `json:"time"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	Url         string `json:"url"`
}

func main() {
	server := http.Server{
		Addr: "localhost:8089",
	}
	http.HandleFunc("/", process)
	server.ListenAndServe()
}

func process(w http.ResponseWriter, r *http.Request) {
	// 1. CREATE CLIENT
	//url := "https://hacker-news.firebaseio.com/v0/topstories.json?print=pretty"
	url := "https://hacker-news.firebaseio.com/v0/item/8863.json?print=pretty"
	client := &http.Client{}

	// 2. CREATE REQUEST
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	// 3. FETCH
	data, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer data.Body.Close()

	// 4. READ BODY (which is io.Reader)
	body, err := ioutil.ReadAll(data.Body)
	if err != nil {
		log.Fatal(err)
	}

	// 5. JSON UNMARSHAL
	article := new(Article)
	err2 := json.Unmarshal(body, article)
	if err2 != nil {
		log.Fatal(err2)
	}

	// MAKE TEMPLATE
	fmt.Printf("%+v", article)
	//
	t, err := template.ParseFiles("tmpl.html")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, article)
}
