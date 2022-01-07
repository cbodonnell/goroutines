package main

import (
	"fmt"
	"time"
)

// main function -- this is a goroutine
func main() {
	// basic_example(5)
	// select_example(25)
	queue_example(20)
}

func basic_example(n int) {
	c := make(chan string)
	go count("sheep", c, n)
	// fmt.Println(<-c) // blocks until a value is read
	for msg := range c {
		fmt.Println(msg)
	}
}

func count(thing string, c chan string, n int) {
	for i := 0; i < n; i++ {
		c <- fmt.Sprintf("%d %s", i, thing)
		time.Sleep(time.Millisecond * 500)
	}
	close(c) // need to close channel to prevent deadlock
}

func select_example(n int) {
	c1 := make(chan string)
	c2 := make(chan string)

	go func() {
		for {
			c1 <- "500 ms"
			time.Sleep(time.Millisecond * 500)
		}
	}()

	go func() {
		for {
			c2 <- "2000 ms"
			time.Sleep(time.Millisecond * 2000)
		}
	}()

	for i := 0; i < n; i++ {
		select { // blocks until a case is ready
		case msg1 := <-c1:
			fmt.Println(msg1)
		case msg2 := <-c2:
			fmt.Println(msg2)
		}
	}
}

func queue_example(n int) {
	jobs := make(chan int, n) // buffered channel
	results := make(chan int, n)

	// spawn a worker
	go worker(jobs, results)

	// add to jobs queue and then close channel
	for i := 0; i < n; i++ {
		jobs <- i
	}
	close(jobs)

	for i := 0; i < n; i++ {
		fmt.Println(<-results)
	}
}

func worker(jobs <-chan int, results chan<- int) {
	for n := range jobs {
		results <- fib(n)
	}
}

func fib(n int) int {
	if n <= 1 {
		return n
	}
	return fib(n-1) + fib(n-2)
}
