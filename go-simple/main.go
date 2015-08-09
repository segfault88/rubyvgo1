package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

const Port = 8000

func root(c web.C, w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("index.html")
	if err != nil {
		panic(err)
	}
	w.Write(data)
}

func slow(c web.C, w http.ResponseWriter, r *http.Request) {
	t1 := time.Now()
	_, err := http.Get(fmt.Sprintf("http://localhost:%d/slow", Port))
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "<span class=\"label label-success\">success %s</span>", time.Now().Sub(t1))
}

func bad(c web.C, w http.ResponseWriter, r *http.Request) {
	t1 := time.Now()
	_, err := http.Get(fmt.Sprintf("http://localhost:%d/bad", Port))
	if err != nil {
		panic(err)
	}
	t := time.Now().Sub(t1)
	if t.Seconds() < 0.5 {
		fmt.Fprintf(w, "<span class=\"label label-success\">good %s</span>", t)
	} else {
		fmt.Fprintf(w, "<span class=\"label label-danger\">bad %s</span>", t)
	}
}

func timeout(c web.C, w http.ResponseWriter, r *http.Request) {
	t1 := time.Now()
	_, err := http.Get(fmt.Sprintf("http://localhost:%d/timeout", Port))
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "<span class=\"label label-danger\">fail %s</span>", time.Now().Sub(t1))
}

func main() {
	goji.Get("/", root)
	goji.Get("/slow", slow)
	goji.Get("/bad", bad)
	goji.Get("/timeout", timeout)
	goji.Serve()
}
