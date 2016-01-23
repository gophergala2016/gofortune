package main

import (
	"fmt"
	_ "html"
	"io/ioutil"
	"net/http"
	"strings"
)

type Handler struct {
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case strings.HasPrefix(r.URL.Path, "/playing-cards/"):
		data, err := ioutil.ReadFile("." + r.URL.Path)
		if err != nil {
			fmt.Fprint(w, err)
		} else {
			w.Header().Add("Content-Type", "image/png")
			w.Write(data)
		}

	case strings.HasPrefix(r.URL.Path, "/js/"):
		data, err := ioutil.ReadFile("." + r.URL.Path)
		if err != nil {
			fmt.Fprint(w, err)
		} else {
			w.Header().Add("Content-Type", "application/javascript")
			w.Write(data)
		}

	case strings.HasPrefix(r.URL.Path, "/css/"):
		data, err := ioutil.ReadFile("." + r.URL.Path)
		if err != nil {
			fmt.Fprint(w, err)
		} else {
			w.Header().Add("Content-Type", "text/css")
			w.Write(data)
		}

	case strings.HasPrefix(r.URL.Path, "/html/"):
		data, err := ioutil.ReadFile("." + r.URL.Path)
		if err != nil {
			fmt.Fprint(w, err)
		} else {
			w.Header().Add("Content-Type", "text/html")
			w.Write(data)
		}

	case r.URL.Path == "/":
		data, err := ioutil.ReadFile("./html/main.html")
		if err != nil {
			fmt.Fprint(w, err)
		} else {
			w.Header().Add("Content-Type", "text/html")
			w.Write(data)
		}
	}
}

func main() {
	var handler Handler

	if err := http.ListenAndServe(":8080", handler); err != nil {
		fmt.Println(err)
	}
}
