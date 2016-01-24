package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Fortune struct {
	wr   http.ResponseWriter
	rq   *http.Request
	deck *Deck
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

	fmt.Println("Listening on http://localhost:8080")

	err := http.ListenAndServe(":8080", &fortune)
	if err != nil {
		fmt.Println(err)
	}
}

func (f *Fortune) ServeHTTP(wr http.ResponseWriter, rq *http.Request) {
	f.wr = wr
	f.rq = rq

	path := rq.URL.Path
	contentType := "text/html"
	root := "."

	switch {
	case strings.HasPrefix(path, "/playing-cards/"):
		contentType = "image/png"

	case strings.HasPrefix(path, "/js/"):
		contentType = "application/javascript"

	case strings.HasPrefix(path, "/css/"):
		contentType = "text/css"

	case strings.HasPrefix(path, "/html/"):
		contentType = "text/html"

	case path == "/":
		contentType = "text/html"
		path = "/html/main.html"

	case path == "/init":
		f.init()
		path = ""

	case path == "/deal":
		f.deal()
		path = ""
	}

	if len(path) > 0 {
		wr.Header().Add("Content-Type", contentType)
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
	f.wr.Header().Add("Content-Type", "application/json")
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
			log.Printf("memorized card: %s\n", response.Card)
		}
	}

	data, err := json.Marshal(response)
	if err != nil {
		response.Error = err.Error()
	}
	f.wr.Header().Add("Content-Type", "application/json")
	f.wr.Write(data)
}
