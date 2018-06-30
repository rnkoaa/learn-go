package main

import (
	"fmt"
	"sync"
	"time"
)

var wg = sync.WaitGroup{}

/* Implementation 1
func main() {
	ch := make(chan int)
	wg.Add(2)

	go func() {
		// receiving data from the channel
		i := <-ch
		fmt.Println(i)
		wg.Done()
	}()

	// Sending the value
	go func() {
		// send go routine
		ch <- 42
		wg.Done()
	}()
	wg.Wait()
}
*/

/* Implementation 2, send multiple items to the channel */
/*
func main() {
	ch := make(chan int)
	for i := 0; i < 5; i++ {
		wg.Add(2)
		go func() {
			// receiving data from the channel
			j := <-ch
			fmt.Println(j)
			wg.Done()
		}()

		// Sending the value
		go func() {
			// send go routine
			ch <- 42
			wg.Done()
		}()
	}
	wg.Wait()
}
*/

/** Implementation 3 Directional Go channel Go Routine */
/*
func main() {
	ch := make(chan int)
	wg.Add(2) // number of go routines

	// receive only channel
	go func(ch <-chan int) {
		i := <-ch
		fmt.Println(i)
		wg.Done()
	}(ch)

	// send only channel
	go func(ch chan<- int) {
		ch <- 42
		wg.Done()
	}(ch)
	wg.Wait()
}
*/

/** Implementation 4 Buffered Go channel Go Routine */
/*
func main() {
	ch := make(chan int, 50)
	wg.Add(2) // number of go routines

	// receive range only channel
	go func(ch <-chan int) {
		for i := range ch {
			fmt.Println(i)
		}
		wg.Done()
	}(ch)

	// send only channel
	go func(ch chan<- int) {
		ch <- 42
		ch <- 27
		close(ch)
		wg.Done()
	}(ch)
	wg.Wait()
}
*/
/** Implementation 5 channel, check response */
/*
func main() {
	ch := make(chan int, 50)
	wg.Add(2) // number of go routines

	// receive range only channel
	go func(ch <-chan int) {
		for {
			if i, ok := <-ch; ok {
				fmt.Println(i)
			} else {
				fmt.Println("Channel Closed.")
				break
			}
		}
		wg.Done()
	}(ch)

	// send only channel
	go func(ch chan<- int) {
		ch <- 42
		ch <- 27
		close(ch)
		wg.Done()
	}(ch)
	wg.Wait()
}
*/

/**
Implementation 6 - Select channels
*/

const (
	logInfo    = "INFO"
	logWarning = "WARNING"
	logERROR   = "ERROR"
)

type logEntry struct {
	time     time.Time
	severity string
	message  string
}

var logCh = make(chan logEntry, 50)
var doneCh = make(chan struct{})

func logger() {
	for {
		select {
		case entry := <-logCh:
			fmt.Printf("%v - [%v] %v\n", entry.time.Format("2006-01-02T15:04:05"),
				entry.severity,
				entry.message)
			break
		case <-doneCh:
			break
		}
	}
}
func main() {
	go logger()
	logCh <- logEntry{time.Now(), logInfo, "App is Starting"}
	logCh <- logEntry{time.Now(), logInfo, "App is Shutting Down"}
	time.Sleep(100 * time.Millisecond)
	doneCh <- struct{}{}
}
