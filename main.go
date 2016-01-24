package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"runtime/debug"
	"sort"
	"strings"
	"time"
)

type Fortune struct {
	wr         http.ResponseWriter
	rq         *http.Request
	deck       *Deck
	scoreCards *Deck
}

func init() {
	debug := flag.Bool("d", false, "debug")
	flag.Parse()

	if !*debug {
		log.SetOutput(ioutil.Discard)
	}
}

func main() {
	var fortune Fortune
	fortune.scoreCards = &Deck{}
	fortune.scoreCards.init()
	fortune.restoreScores()

	fmt.Println("Listening on http://localhost:8080")

	err := http.ListenAndServe(":8080", &fortune)
	if err != nil {
		fmt.Println(err)
	}
}

func (f *Fortune) ServeHTTP(wr http.ResponseWriter, rq *http.Request) {
	defer func() {
		obj := recover()
		if obj != nil {
			msg := fmt.Sprintf("<pre>Error: %v\nStack: %v</pre>", obj, string(debug.Stack()))
			io.WriteString(wr, msg)
			fmt.Println(msg)
		}
	}()

	f.wr = wr
	f.rq = rq

	path := rq.URL.Path
	contentType := "text/html"
	root := "."

	switch {
	case strings.HasPrefix(path, "/playing-cards/"):
		contentType = "image/png"
		wr.Header().Set("cache-control", "max-age=86400")

	case strings.HasPrefix(path, "/js/"):
		contentType = "application/javascript"

	case strings.HasPrefix(path, "/css/"):
		contentType = "text/css"

	case strings.HasPrefix(path, "/html/"):
		contentType = "text/html"

	case path == "/":
		contentType = "text/html"
		path = "/html/main.html"
		fmt.Printf("%s: Visitor from %s\n", time.Now(), rq.RemoteAddr)

	case path == "/init":
		f.init()
		path = ""

	case path == "/deal":
		f.deal()
		path = ""

	case path == "/fortune":
		f.fortune()
		path = ""

	case path == "/scores":
		f.scores()
		path = ""
	}

	if len(path) > 0 {
		wr.Header().Set("Content-Type", contentType)
		data, err := ioutil.ReadFile(root + path)
		if err == nil {
			wr.Write(data)
		} else {
			fmt.Fprint(wr, err)
		}
	}
}

func (f *Fortune) init() {
	f.deck = &Deck{}
	f.deck.init()
	f.deck.shuffle()
	f.deck.Cards = f.deck.Cards[:21]

	type Response struct {
		Cards []*Card
		Error string
	}

	response := &Response{
		Cards: f.deck.Cards,
	}
	data, err := json.Marshal(response)
	if err != nil {
		response.Error = err.Error()
	}
	f.wr.Header().Set("Content-Type", "application/json")
	f.wr.Write(data)
}

func (f *Fortune) deal() {
	type RequestCard struct {
		Image string
	}
	type Request struct {
		Cards []RequestCard
		Row   int
		Count int
	}
	type Response struct {
		Row1  []*Card
		Row2  []*Card
		Row3  []*Card
		Card  string
		Error string
	}

	response := &Response{}
	reqData, err := ioutil.ReadAll(f.rq.Body)
	if err != nil {
		response.Error = err.Error()
	} else {
		request := &Request{}
		err = json.Unmarshal(reqData, request)
		if err != nil {
			response.Error = err.Error()
		}
		f.deck = &Deck{}
		for _, card := range request.Cards {
			f.deck.Cards = append(f.deck.Cards, &Card{Image: card.Image})
		}
		if len(request.Cards) == 21 {
			if request.Row == 0 {
				response.Row1 = f.deck.Cards[:7]
				response.Row2 = f.deck.Cards[7:14]
				response.Row3 = f.deck.Cards[14:]
			} else {
				f.deck.placeMiddle(request.Row)
				f.deck.deal()
				response.Row1 = f.deck.Row1
				response.Row2 = f.deck.Row2
				response.Row3 = f.deck.Row3
			}
		} else {
			response.Error += "\nDeck should have 21 cards."
		}
		log.Printf("request: %v\n", request)
		if request.Count == 3 {
			response.Card = f.deck.Row2[3].Image
			f.scoreCard(response.Card)
			log.Printf("memorized card: %s\n", response.Card)
		}
	}

	data, err := json.Marshal(response)
	if err != nil {
		response.Error = err.Error()
	}
	f.wr.Header().Set("Content-Type", "application/json")
	f.wr.Write(data)
}

func (f *Fortune) fortune() {
	words := map[string]string{
		"2C.png":  "passion",
		"2D.png":  "wealth",
		"2H.png":  "love",
		"2S.png":  "law",
		"3C.png":  "interest",
		"3D.png":  "rich",
		"3H.png":  "like",
		"3S.png":  "rule",
		"4C.png":  "positive",
		"4D.png":  "gold",
		"4H.png":  "nice",
		"4S.png":  "command",
		"5C.png":  "real",
		"5D.png":  "money",
		"5H.png":  "related",
		"5S.png":  "advise",
		"6C.png":  "growing",
		"6D.png":  "fortune",
		"6H.png":  "good",
		"6S.png":  "statement",
		"7C.png":  "study",
		"7D.png":  "well",
		"7H.png":  "sweet",
		"7S.png":  "court",
		"8C.png":  "understand",
		"8D.png":  "cash",
		"8H.png":  "protect",
		"8S.png":  "action",
		"9C.png":  "hobby",
		"9D.png":  "stock",
		"9H.png":  "live",
		"9S.png":  "act",
		"10C.png": "knowledge",
		"10D.png": "value",
		"10H.png": "friend",
		"10S.png": "order",
		"JC.png":  "student",
		"JD.png":  "banker",
		"JH.png":  "husband",
		"JS.png":  "judge",
		"QC.png":  "nurse",
		"QD.png":  "actress",
		"QH.png":  "wife",
		"QS.png":  "queen",
		"KC.png":  "researcher",
		"KD.png":  "ceo",
		"KH.png":  "lover",
		"KS.png":  "congressman",
		"AC.png":  "president",
		"AD.png":  "thesaurus",
		"AH.png":  "family",
		"AS.png":  "country",
	}
	type Request struct {
		Card string
	}
	type Response struct {
		Tweet string
		Error string
	}

	response := &Response{}

	request := &Request{}
	reqData, err := ioutil.ReadAll(f.rq.Body)
	err = json.Unmarshal(reqData, request)
	if err != nil {
		response.Error = err.Error()
	}

	key, _ := words[request.Card]
	url := "https://twitter.com/search?f=realtime&q=" + key
	resp, err := http.Get(url)
	if err != nil {
		response.Error = "Error: " + err.Error()
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		response.Error = "Error: " + err.Error()
	}
	resp.Body.Close()

	search := regexp.MustCompile(`<p class="TweetTextSize .*>.*</p>`)
	tweets := search.FindAllString(string(body), -1)

	tweet := "Unable to fetch tweets."
	if len(tweets) > 0 {
		index := rand.Intn(len(tweets))
		tweet = tweets[index]
	}
	response.Tweet = tweet
	fmt.Printf("Visitor=%s word=%s fortune=%s\n", f.rq.RemoteAddr, key, tweet)

	data, err := json.Marshal(response)
	if err != nil {
		response.Error = err.Error()
	}
	f.wr.Header().Set("Content-Type", "application/json")
	f.wr.Write(data)
}

func (f *Fortune) scoreCard(memorizedCard string) {
	for _, card := range f.scoreCards.Cards {
		if card.Image == memorizedCard {
			card.Score++
			break
		}
	}
	f.saveScores()
}

func (f *Fortune) scores() {
	type Response struct {
		ScoreCards []*Card
		Error      string
	}

	response := &Response{}
	sort.Sort(f.scoreCards)
	response.ScoreCards = f.scoreCards.Cards

	data, err := json.Marshal(response)
	if err != nil {
		response.Error = err.Error()
	}
	f.wr.Header().Set("Content-Type", "application/json")
	f.wr.Write(data)
}

func (f *Fortune) saveScores() {
	data, err := json.Marshal(f.scoreCards)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ioutil.WriteFile("scores.json", data, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (f *Fortune) restoreScores() {
	data, err := ioutil.ReadFile("scores.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(data, f.scoreCards)
	if err != nil {
		fmt.Println(err)
		return
	}
}
