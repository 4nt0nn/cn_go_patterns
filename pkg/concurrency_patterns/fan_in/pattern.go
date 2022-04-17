package fan_in

import (
	"sync"
)

// Funnel starts a dedicated goroutine for each channel in list of sources
// that reads values from its assigned channel and forwards them to dest, a single-output
// channel shared by all of the goroutine
func Funnel(sources ...<-chan int) <-chan int {
	dest := make(chan int) // The shared output channel

	var wg sync.WaitGroup // Used to automatically close dest when all sources are closed.

	wg.Add(len(sources)) // Set size of the WaitGroup

	for _, ch := range sources { // Start a goroutine for each source
		go func(c <-chan int) {
			defer wg.Done() // Notify WaitGroup when c closes.

			for n := range c {
				dest <- n
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(dest)
	}()

	return dest
}

////// Example impl /////

// func main() {
// 	sources := make([]<-chan int, 0) // Create an empty channel slice

// 	for i := 0; i < 3; i++ {
// 		ch := make(chan int)
// 		sources = append(sources, ch) // Create a channel; add to sources

// 		go func() { // Run a toy goroutine for each
// 			defer close(ch) // Close ch when the routine ends

// 			for i := 1; i <= 5; i++ {
// 				ch <- i
// 				time.Sleep(time.Second)
// 			}
// 		}()
// 	}

// 	dest := Funnel(sources...)
// 	for d := range dest {
// 		fmt.Println(d)
// 	}
// }
