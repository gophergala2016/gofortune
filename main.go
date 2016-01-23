package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Fortune struct {
	wr   http.ResponseWriter
	rq   *http.Request
	deck *Deck
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
		wr.Header().Add("Content-Type", "application/json")
		f.init()
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
	type Response struct {
		Cards []*Card
		Error string
	}
	f.deck.Cards = f.deck.Cards[:21]
	response := &Response{
		Cards: f.deck.Cards,
	}
	data, err := json.Marshal(response)
	if err != nil {
		response.Error = err.Error()
	}
	f.wr.Write(data)
}

func main() {
	var fortune Fortune

	fmt.Println("Listening on http://localhost:8080")

	err := http.ListenAndServe(":8080", &fortune)
	if err != nil {
		fmt.Println(err)
	}
}
