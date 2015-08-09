package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

const APIPort = 8000

type Scanner struct {
	results chan Result
	done    bool
}

type Result struct {
	Service int    `json:"service"`
	Api     int    `json:"api"`
	Css     string `json:"css"`
	Text    string `json:"text"`
}

var (
	addr     = flag.String("addr", ":8888", "http service address")
	scanners = make(map[int]*Scanner)
	lock     sync.Mutex
)

func rootRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data, err := ioutil.ReadFile("index.html")
	panicIfErr(err)
	w.Write(data)
}

func scanRoute(w http.ResponseWriter, r *http.Request) {
	key := int(rand.Int31())
	// send the key for this scan
	fmt.Fprintf(w, "%d", key)
	log.Printf("Starting scan, key: %d", key)

	scanner := new(Scanner)
	scanner.results = make(chan Result, 16)

	// store the scanner in the map to be looked up by the poll route
	func() {
		lock.Lock()
		defer lock.Unlock()
		scanners[key] = scanner
	}()

	go scanServices(scanner)
}

func pollRoute(w http.ResponseWriter, r *http.Request) {
	key, err := strconv.Atoi(mux.Vars(r)["key"])
	panicIfErr(err)

	// look up the scanner
	scanner := func() *Scanner {
		lock.Lock()
		defer lock.Unlock()
		return scanners[key]
	}()

	if scanner == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response := struct {
		Messages []Result `json:"results"`
		Done     bool     `json:"done"`
	}{
		make([]Result, 0),
		scanner.done,
	}

	// grab all the waiting results
	func() {
		for {
			select {
			case m := <-scanner.results:
				response.Messages = append(response.Messages, m)
			default:
				return
			}
		}
	}()

	// marshal and send json
	json, err := json.Marshal(&response)
	panicIfErr(err)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
	log.Printf("Poll result key: %d, count: %d, done: %t\n", key, len(response.Messages), response.Done)
}

func scanServices(scanner *Scanner) {
	var wg sync.WaitGroup

	for i := 1; i < 11; i++ {
		for j := 1; j < 6; j++ {
			wg.Add(1)
			if j == 1 || j == 2 {
				go callSlow(i, j, &wg, &scanner.results)
			} else if j == 3 || j == 4 {
				go callBad(i, j, &wg, &scanner.results)
			} else {
				go callTimeout(i, j, &wg, &scanner.results)
			}
		}
	}

	wg.Wait()
	scanner.done = true
	log.Println("Scan done")
}

func callSlow(i int, j int, wg *sync.WaitGroup, r *chan Result) {
	defer wg.Done()
	t1 := time.Now()
	_, err := http.Get(fmt.Sprintf("http://localhost:%d/slow", APIPort))
	panicIfErr(err)
	*r <- Result{i, j, "success", fmt.Sprintf("success %s", time.Now().Sub(t1))}
}

func callBad(i int, j int, wg *sync.WaitGroup, r *chan Result) {
	defer wg.Done()
	t1 := time.Now()
	_, err := http.Get(fmt.Sprintf("http://localhost:%d/bad", APIPort))
	panicIfErr(err)
	t := time.Now().Sub(t1)
	if t.Seconds() < 0.5 {
		*r <- Result{i, j, "success", fmt.Sprintf("good %s", t)}
	} else {
		*r <- Result{i, j, "warning", fmt.Sprintf("bad %s", t)}
	}
}

func callTimeout(i int, j int, wg *sync.WaitGroup, r *chan Result) {
	defer wg.Done()
	t1 := time.Now()
	_, err := http.Get(fmt.Sprintf("http://localhost:%d/timeout", APIPort))
	panicIfErr(err)
	*r <- Result{i, j, "danger", fmt.Sprintf("fail %s", time.Now().Sub(t1))}
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/", rootRoute)
	r.HandleFunc("/scan", scanRoute)
	r.HandleFunc("/poll/{key}", pollRoute)

	http.Handle("/", r)

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}
