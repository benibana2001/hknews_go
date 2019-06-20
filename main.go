package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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
	urlTopStories := "https://hacker-news.firebaseio.com/v0/topstories.json?print=pretty"
	//urlItem := "https://hacker-news.firebaseio.com/v0/item/8863.json?print=pretty"
	urlItemBase := "https://hacker-news.firebaseio.com/v0/item/"
	client := &http.Client{}

	// 2a. CREATE REQUEST
	reqTopStories, err := http.NewRequest("GET", urlTopStories, nil)
	if err != nil {
		log.Fatal(err)
	}

	// 3a. FETCH
	dataTopStories, err := client.Do(reqTopStories)
	if err != nil {
		log.Fatal(err)
	}

	defer dataTopStories.Body.Close()

	// 4a. READ BODY (which is io.Reader)
	bodyTopStories, err := ioutil.ReadAll(dataTopStories.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(bodyTopStories)

	var ids []int
	json.Unmarshal(bodyTopStories, &ids)
	fmt.Printf("%+v", ids)

	// 2b. CREATE REQUEST
	//reqItem, err := http.NewRequest("GET", urlItem, nil)
	id := strconv.Itoa(ids[0])
	reqItem, err := http.NewRequest("GET", urlItemBase+id+".json", nil)
	if err != nil {
		log.Fatal(err)
	}

	// 3b. FETCH
	dataItem, err := client.Do(reqItem)
	if err != nil {
		log.Fatal(err)
	}

	defer dataItem.Body.Close()

	// 4b. READ BODY (which is io.Reader)
	bodyItem, err := ioutil.ReadAll(dataItem.Body)
	if err != nil {
		log.Fatal(err)
	}

	// 5b. JSON UNMARSHAL
	article := new(Article)
	err2 := json.Unmarshal(bodyItem, article)
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
