package Controllers

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

const (
	urlTopStories = "https://hacker-news.firebaseio.com/v0/topstories.json?print=pretty"
	urlItemBase   = "https://hacker-news.firebaseio.com/v0/item/"
)
