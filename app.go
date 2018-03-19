package main

import (
	"time"
	"math/rand"
	"sync"
	"sort"
	"net/http"
	"html/template"
	"log"
	"github.com/gorilla/mux"
	"encoding/json"
)

const MaxPairsCapacity = 20
const DiscardAmountForGC = 6
const LatestPairsDisplayAmount = 5
const OldestSumAmount = 4

type Pair struct {
	Second int
	Value  int
}

var pairs []Pair
var sumList []int
var wg sync.WaitGroup
var mu sync.Mutex

func GenPair(ch chan Pair) {
	var p Pair
	for {
		p = Pair{rand.Intn(100-1) + 1, rand.Intn(10000-1000) + 1000}
		ch <- p
		time.Sleep(time.Second)
	}
}

func Process(ch chan Pair) {
	for {
		x := <-ch
		mu.Lock()
		pairs = append(pairs, x)
		mu.Unlock()
		if len(pairs) == MaxPairsCapacity {
			wg.Add(1)
			go ForceDiscard(DiscardAmountForGC)
		}
		wg.Add(1)
		go DiscardAfter(x)
	}
}

func DiscardAfter(x Pair) {
	time.Sleep(time.Second * time.Duration(x.Second))
	var na []Pair
	mu.Lock()
	for _, v := range pairs {
		if v == x {
			continue
		} else {
			na = append(na, v)
		}
	}
	pairs = na
	mu.Unlock()
}

func GetOldest(amount int) (l []Pair) {
	l = pairs[:amount]
	return
}

func GetLatest(amount int) (l []Pair) {
	l = pairs[len(pairs)-amount:len(pairs)]
	return
}

func ForceDiscard(amount int) {
	mu.Lock()
	pairs = GetLatest(len(pairs) - amount)
	mu.Unlock()
}

func SumOldest(amount int) (total int) {
	mu.Lock()
	if len(pairs) >= amount {
		for _, pair := range GetOldest(amount) {
			total += pair.Value
		}

	}
	sumList = append(sumList, total)
	mu.Unlock()
	return
}

func GetMedian(l []int) (median int) {
	if len(l) == 0 {
		return
	}

	sort.Ints(l)
	middle := len(l) / 2
	median = l[middle]
	if len(l)%2 == 0 {
		median = (median + l[middle-1]) / 2
	}
	return
}

// ---------------------------- REST APIs -------------------------------

func LoadMainPage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("assets/main.html")
	t.Execute(w, nil)
}

func HandleGetLatestPairs(w http.ResponseWriter, r *http.Request) {
	amount := LatestPairsDisplayAmount
		if len(pairs) >= amount {
			json.NewEncoder(w).Encode(GetLatest(amount))
		} else {
			json.NewEncoder(w).Encode(pairs)
		}
}

func HandleGetSum(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(SumOldest(OldestSumAmount))
}

func HandleGetMedian(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(GetMedian(sumList))
}

// ------------------------------ MAIN ---------------------------------

func main() {
	wg.Add(6)

	go func() {
		router := mux.NewRouter()
		router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))
		router.HandleFunc("/", LoadMainPage)
		router.HandleFunc("/pairs", HandleGetLatestPairs).Methods("GET")
		router.HandleFunc("/sum", HandleGetSum).Methods("GET")
		router.HandleFunc("/median", HandleGetMedian).Methods("GET")

		log.Fatal(http.ListenAndServe(":8080", router))
	}()

	c := make(chan Pair)

	go GenPair(c)
	go Process(c)

	wg.Wait()
}
