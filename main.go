package main

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
type Ids []int // Article に Cards を格納する際の中間スライス
type Cards []Article

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	server := http.Server{
		Addr: ":" + port,
	}
	http.HandleFunc("/", process)
	server.ListenAndServe()
}

func process(w http.ResponseWriter, r *http.Request) {
	// 1. CREATE CLIENT
	urlTopStories := "https://hacker-news.firebaseio.com/v0/topstories.json?print=pretty"
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
	bodyTopStories, err := ioutil.ReadAll(dataTopStories.Body) // bodyTopStories is `[]bytes`
	if err != nil {
		log.Fatal(err)
	}

	var ids Ids
	json.Unmarshal(bodyTopStories, &ids)
	ids = ids[0:9]

	// 2c. CREATE REQUEST , 3c. FETCH , 4c. READ BODY (which is io.Reader), 5c. JSON UNMARSHAL
	cards := Cards{}

	for _, id := range ids {
		reqItem, err := http.NewRequest("GET", urlItemBase+strconv.Itoa(id)+".json", nil)
		if err != nil {
			log.Fatal(err)
		}

		dataItem, err := client.Do(reqItem)
		if err != nil {
			log.Fatal(err)
		}

		storeArticle(dataItem, &cards) // cards に個々の記事を格納する
	}

	// 6. MAKE TEMPLATE
	t, err := template.ParseFiles("./template/tmpl.html")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, cards)
}

func storeArticle(item *http.Response, cards *Cards) {
	defer item.Body.Close()

	bodyItem, err := ioutil.ReadAll(item.Body)
	if err != nil {
		log.Fatal(err)
	}
	err1 := item.Body.Close()
	if err1 != nil {
		log.Fatal(err)
	}

	article := new(Article)
	err2 := json.Unmarshal(bodyItem, article)
	if err != nil {
		log.Fatal(err2)
	}
	*cards = append(*cards, *article)
}
