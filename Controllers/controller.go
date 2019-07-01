package Controllers

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Controller struct{}

// Home is http handler and use html/template
func (c Controller) Home(w http.ResponseWriter, r *http.Request) {
	reader := getBody()
	defer reader.Close()
	cards := parseToCards(reader)
	fmt.Printf("type of cards : %T", cards)

	// MAKE TEMPLATE
	t, err := template.ParseFiles("./template/tmpl.html")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, cards)
}

func getBody() io.ReadCloser {
	// 1. CREATE CLIENT
	client := &http.Client{}

	// 2. CREATE REQUEST
	reqTopStories, err := http.NewRequest("GET", urlTopStories, nil)
	if err != nil {
		log.Fatal(err)
	}

	// 3. FETCH
	dataTopStories, err := client.Do(reqTopStories)
	if err != nil {
		log.Fatal(err)
	}

	return dataTopStories.Body
}

func fetchArticleId(r io.ReadCloser) Ids {
	bodyTopStories, err := ioutil.ReadAll(r) // bodyTopStories is `[]bytes`
	if err != nil {
		log.Fatal(err)
	}
	var ids Ids
	json.Unmarshal(bodyTopStories, &ids)
	ids = ids[0:9] // ids == []int
	return ids
}

func parseToCards(r io.ReadCloser) *Cards {
	ids := fetchArticleId(r)
	client := &http.Client{}

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

		storeArticle(dataItem, &cards)
	}

	return &cards
}

// storeArticle は cards に個々の記事を格納します。
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
