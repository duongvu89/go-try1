package main

import (
	"fmt"
	"time"
	"math/rand"
	"sync"
	"sort"
)

type Pair struct {
	first int
	second int
}

var a []Pair
var sumList []int
var wg sync.WaitGroup
var mu sync.Mutex

func genPair(ch chan Pair) {
	var p Pair
	for {
		p = Pair{rand.Intn(50 - 1) + 1, rand.Intn(10000 - 1000) + 1000}
		ch <- p
		time.Sleep(time.Second)
	}
}

func appendList(ch chan Pair, quit chan bool, gc chan bool) {
	for {
		select {
		case x := <- ch:
			mu.Lock()
			a = append(a, x)
			mu.Unlock()
			go discardAfter(x)
		case <- quit:
			fmt.Println(a)
			return
		case <- gc:
			go forceDiscard(3)
		}

	}
}

func discardAfter(x Pair) {
	// na = new a, x = a that's to be deleted
	time.Sleep(time.Second * time.Duration(x.first))
	var na []Pair
	mu.Lock()
	for _, v := range a {
		if v == x {
			continue
		} else {
			na = append(na, v)
		}
	}
	a = na
	mu.Unlock()
}

func getOldest(amount int) (l []Pair) {
	l = a[:amount]
	return
}

func getLatest(amount int) (l []Pair) {
	l = a[len(a) - amount:len(a)]
	return
}

func forceDiscard(amount int) {
	mu.Lock()
	a = getLatest(len(a) - amount)
	mu.Unlock()
}

func sumOldest(amount int) (total int) {
	mu.Lock()
	if len(a) >= 2 {
		fmt.Printf("The latest array %v with ", a)
		for _, pair := range getOldest(amount) {
			total += pair.second
		}

	}
	sumList = append(sumList, total)
	mu.Unlock()
	return
}

func pressSumButton() {
	go func() {
		for {
			time.Sleep(time.Second * 2)
			fmt.Println("sum=", sumOldest(2))
		}
	}()
}

func getMedian(l []int) (median int) {
	sort.Ints(l)
	middle := len(l) / 2
	result := l[middle]
	if len(l)%2 == 0 {
		result = (result + l[middle-1]) / 2
	}
	return result
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
	quit := make(chan bool)
	gc := make (chan bool)

	go genPair(c)
	go func() {
		for {
			time.Sleep(time.Second / 2)
			fmt.Println(a)
			if len(a) == 7 {
				quit <- true
			}
			if len(a) == 5 {
				gc <- true
			}
		}
	}()

	pressSumButton()
	pressMedianButton()

	appendList(c, quit, gc)
}