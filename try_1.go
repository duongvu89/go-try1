package main

import (
	"fmt"
	"time"
	"math/rand"
)

type Pair struct {
	first int
	second int
}

var a []Pair


func genPair(ch chan Pair) {
	var p Pair
	for {
		p = Pair{rand.Intn(10 - 1) + 1, rand.Intn(10000 - 1000) + 1000}
		ch <- p
		time.Sleep(time.Second)
	}
}

func appendList(ch chan Pair, quit chan bool) {
	for {
		select {
		case x := <- ch:
			a = append(a, x)
			go discardAfter(x, a)

		case <- quit:
			fmt.Println(a)
			return
		}

	}
}

func discardAfter(x Pair, a [] Pair) {
	// na = new a, x = a that's to be deleted
	time.Sleep(time.Second * time.Duration(x.first))
	var na []Pair
	for _, v := range a {
		if v == x {
			continue
		} else {
			na = append(na, v)
		}
	}
	a = na
}


func main() {
	c := make(chan Pair)
	quit := make(chan bool)

	go genPair(c)
	go func() {
		time.Sleep(10 * time.Second)
		quit <- true
	}()

	appendList(c, quit)
}