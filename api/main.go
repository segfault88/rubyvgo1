package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

func slow(c web.C, w http.ResponseWriter, r *http.Request) {
	n := time.Duration(200 + (rand.Int31n(200)))
	time.Sleep(n * time.Millisecond)
	fmt.Fprintf(w, "Done %d ms", n)
}

func bad(c web.C, w http.ResponseWriter, r *http.Request) {
	if rand.Float32() <= 0.5 {
		n := time.Duration(50 + (rand.Int31n(50)))
		time.Sleep(n * time.Millisecond)
		fmt.Fprintf(w, "Done %d ms", n)
	} else {
		time.Sleep(3 * time.Second)
		fmt.Fprint(w, "This api is bad!")
	}
}

func timeout(c web.C, w http.ResponseWriter, r *http.Request) {
	time.Sleep(10 * time.Second)
	fmt.Fprint(w, "That took a while!")
}

func main() {
	rand.Seed(time.Now().UnixNano())
	goji.Get("/slow", slow)
	goji.Get("/bad", bad)
	goji.Get("/timeout", timeout)
	goji.Serve()
}
