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
type Ids []int // Article に Cards を格納する際の中間スライス
type Cards []Article

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
	fmt.Printf("%+v\n", ids)

	// 2b. CREATE REQUEST
	//reqItem, err := http.NewRequest("GET", urlItem, nil)
	id := strconv.Itoa(ids[0])
	reqItem, err := http.NewRequest("GET", urlItemBase+id+".json", nil)
	if err != nil {
		log.Fatal(err)
	}

	//////////
	//////////

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

		// TODO: SHOULD NOT CALL DEFER IN LOOP, REF: https://blog.golang.org/defer-panic-and-recover
		//  => 無名関数でラップする
		//defer dataItem.Body.Close()

		bodyItem, err := ioutil.ReadAll(dataItem.Body)
		if err != nil {
			log.Fatal(err)
		}
		// TODO: DEFER
		err1 := dataItem.Body.Close()
		if err1 != nil {
			log.Fatal(err)
		}

		article := new(Article)
		err2 := json.Unmarshal(bodyItem, article)
		if err != nil {
			log.Fatal(err2)
		}

		cards = append(cards, *article)
	}
	fmt.Printf("%+v\n", cards)
	//////////
	//////////

	//////////
	//////////
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
	//////////
	//////////

	// MAKE TEMPLATE
	fmt.Printf("%+v\n", article)
	//
	t, err := template.ParseFiles("tmpl.html")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, article)
}
