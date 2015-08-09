package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const APIPort = 8000

type Result struct {
	Service int    `json:"service"`
	Api     int    `json:"api"`
	Css     string `json:"css"`
	Text    string `json:"text"`
}

var (
	addr     = flag.String("addr", ":8889", "http service address")
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func rootRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data, err := ioutil.ReadFile("index.html")
	panicIfErr(err)
	w.Write(data)
}

func wsRoute(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	scanServices(ws)
}

func scanServices(ws *websocket.Conn) {
	defer ws.Close()

	messages := make(chan Result, 16)
	var wg sync.WaitGroup

	for i := 1; i < 11; i++ {
		for j := 1; j < 6; j++ {
			wg.Add(1)
			if j == 1 || j == 2 {
				go callSlow(i, j, &wg, messages)
			} else if j == 3 || j == 4 {
				go callBad(i, j, &wg, messages)
			} else {
				go callTimeout(i, j, &wg, messages)
			}
		}
	}

	done := make(chan bool)
	go func() {
		wg.Wait()
		done <- true
	}()

	for {
		select {
		case msg := <-messages:
			json, err := json.Marshal(msg)
			panicIfErr(err)
			ws.WriteMessage(websocket.TextMessage, json)
		case <-done:
			log.Println("Scan done")
			break
		}
	}
}

func callSlow(i int, j int, wg *sync.WaitGroup, r chan Result) {
	defer wg.Done()
	t1 := time.Now()
	_, err := http.Get(fmt.Sprintf("http://localhost:%d/slow", APIPort))
	panicIfErr(err)
	r <- Result{i, j, "success", fmt.Sprintf("success %s", time.Now().Sub(t1))}
}

func callBad(i int, j int, wg *sync.WaitGroup, r chan Result) {
	defer wg.Done()
	t1 := time.Now()
	_, err := http.Get(fmt.Sprintf("http://localhost:%d/bad", APIPort))
	panicIfErr(err)
	t := time.Now().Sub(t1)
	if t.Seconds() < 0.5 {
		r <- Result{i, j, "success", fmt.Sprintf("good %s", t)}
	} else {
		r <- Result{i, j, "warning", fmt.Sprintf("bad %s", t)}
	}
}

func callTimeout(i int, j int, wg *sync.WaitGroup, r chan Result) {
	defer wg.Done()
	t1 := time.Now()
	_, err := http.Get(fmt.Sprintf("http://localhost:%d/timeout", APIPort))
	panicIfErr(err)
	r <- Result{i, j, "danger", fmt.Sprintf("fail %s", time.Now().Sub(t1))}
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()

	http.HandleFunc("/", rootRoute)
	http.HandleFunc("/ws", wsRoute)

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}
