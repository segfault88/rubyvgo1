package main

import (
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

var (
	addr     = flag.String("addr", ":8880", "http service address")
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

	messages := make(chan string, 16)
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
			ws.WriteMessage(websocket.TextMessage, []byte(msg))
		case <-done:
			log.Println("Scan done")
			break
		}
	}
}

func callSlow(i int, j int, wg *sync.WaitGroup, r chan string) {
	defer wg.Done()
	t1 := time.Now()
	_, err := http.Get(fmt.Sprintf("http://localhost:%d/slow", APIPort))
	panicIfErr(err)
	r <- fmtResut(i, j, "info", "success", time.Now().Sub(t1))
}

func callBad(i int, j int, wg *sync.WaitGroup, r chan string) {
	defer wg.Done()
	t1 := time.Now()
	_, err := http.Get(fmt.Sprintf("http://localhost:%d/bad", APIPort))
	panicIfErr(err)
	t := time.Now().Sub(t1)
	if t.Seconds() < 0.5 {
		r <- fmtResut(i, j, "success", "good", t)
	} else {
		r <- fmtResut(i, j, "warning", "bad", t)
	}
}

func callTimeout(i int, j int, wg *sync.WaitGroup, r chan string) {
	defer wg.Done()
	t1 := time.Now()
	_, err := http.Get(fmt.Sprintf("http://localhost:%d/timeout", APIPort))
	panicIfErr(err)
	r <- fmtResut(i, j, "danger", "fail", time.Now().Sub(t1))
}

func fmtResut(i int, j int, class string, msg string, t time.Duration) string {
	return fmt.Sprintf("$('.s%d .a%d').html('<span class=\"label label-%s\">%s %s</span>')", i, j, class, msg, t)
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
