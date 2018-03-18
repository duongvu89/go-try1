package main

import (
	"fmt"
	"time"
	"math/rand"
	"sync"
	"sort"
)

const MaxPairsCapacity = 20
const DiscardAmountForGC = 6
const LatestPairsDisplayAmount = 5
const OldestSumAmount = 4

type Pair struct {
	second int
	value  int
}

var pairs []Pair
var sumList []int
var wg sync.WaitGroup
var mu sync.Mutex

func genPair(ch chan Pair) {
	var p Pair
	for {
		p = Pair{rand.Intn(100-1) + 1, rand.Intn(10000-1000) + 1000}
		ch <- p
		time.Sleep(time.Second)
	}
}

func daemon(ch chan Pair) {
	for {
		x := <-ch
		mu.Lock()
		pairs = append(pairs, x)
		mu.Unlock()
		if len(pairs) == MaxPairsCapacity {
			wg.Add(1)
			go forceDiscard(DiscardAmountForGC)
		}
		wg.Add(1)
		go discardAfter(x)
	}
}

func discardAfter(x Pair) {
	time.Sleep(time.Second * time.Duration(x.second))
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

func getOldest(amount int) (l []Pair) {
	l = pairs[:amount]
	return
}

func getLatest(amount int) (l []Pair) {
	l = pairs[len(pairs)-amount:len(pairs)]
	return
}

func displayLatestPair(amount int) {
	for {
		if len(pairs) >= amount {
			fmt.Printf("Latest %d pairs: %v \n", amount, getLatest(amount))
		} else {
			fmt.Printf("Latest %d pairs: %v \n", len(pairs), pairs)
		}
		time.Sleep(time.Second)
	}
}

func forceDiscard(amount int) {
	mu.Lock()
	pairs = getLatest(len(pairs) - amount)
	mu.Unlock()
}

func sumOldest(amount int) (total int) {
	mu.Lock()
	if len(pairs) >= amount {
		for _, pair := range getOldest(amount) {
			total += pair.value
		}

	}
	sumList = append(sumList, total)
	mu.Unlock()
	return
}

func pressSumButton() {
	go func() {
		for {
			time.Sleep(time.Second * 5)
			fmt.Println("sum=", sumOldest(OldestSumAmount))
		}
	}()
}

func getMedian(l []int) (median int) {
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

func pressMedianButton() {
	go func() {
		for {
			time.Sleep(time.Second * 3)
			fmt.Printf("sumList %v median %d \n", sumList, getMedian(sumList))
		}
	}()
}

func main() {
	c := make(chan Pair)

	wg.Add(5)
	go genPair(c)

	go daemon(c)

	go displayLatestPair(LatestPairsDisplayAmount)

	pressSumButton()
	pressMedianButton()

	wg.Wait()
}
