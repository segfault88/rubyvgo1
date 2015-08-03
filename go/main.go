package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const Port = 8000

var (
	addr = flag.String("addr", ":8880", "http service address")
)

func root(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("index.html")
	if err != nil {
		panic(err)
	}
	w.Write(data)
}

func slow(w http.ResponseWriter, r *http.Request) {
	t1 := time.Now()
	_, err := http.Get(fmt.Sprintf("http://localhost:%d/slow", Port))
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "<span class=\"label label-success\">success %s</span>", time.Now().Sub(t1))
}

func bad(w http.ResponseWriter, r *http.Request) {
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

func timeout(w http.ResponseWriter, r *http.Request) {
	t1 := time.Now()
	_, err := http.Get(fmt.Sprintf("http://localhost:%d/timeout", Port))
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "<span class=\"label label-danger\">fail %s</span>", time.Now().Sub(t1))
}

func main() {
	flag.Parse()

	http.HandleFunc("/", root)
	http.HandleFunc("/slow", slow)
	http.HandleFunc("/bad", bad)
	http.HandleFunc("/timeout", timeout)

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}
