package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Write 3 functions:
// writer - generates numbers 1 to 10
// doubler - multiplies by 2, imitating work (500ms)
// reader - reads and prints to the console

func writer() <-chan int {
	ch := make(chan int)

	go func() {
		for range 10 {
			randNum := 1 + rand.Intn(10)
			ch <- randNum
		}
		close(ch)
	}()

	return ch
}

func doubler(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		for v := range in {
			time.Sleep(500 * time.Millisecond)
			out <- v * 2
		}
		close(out)
	}()

	return out
}

func reader(in <-chan int) {
	for v := range in {
		fmt.Println("v =", v)
	}
}

func main() {
	pipeline := doubler(writer())
	reader(pipeline)
}
