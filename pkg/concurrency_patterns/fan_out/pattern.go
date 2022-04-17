package fan_out

import (
	"fmt"
	"sync"
)

// Split creates the destination channels. For each channel created, it executes a goroutine that
// retrieves values from source in a for loop and forwards them to their assigned output channel.
// Effectively, each goroutine competes for reads from source; if several are trying to read, the "winner"
// will be randomly determined. If source is closed, all goroutines terminate and all of the destination channels are closed.
func Split(source <-chan int, n int) []<-chan int {
	dests := make([]<-chan int, 0) // Create the dests slice

	for i := 0; i < n; i++ { // Create n destination channles.
		ch := make(chan int)
		dests = append(dests, ch)

		go func() { // Each channel gets a dedicated goroutine that competes for reads
			defer close(ch)

			for val := range source {
				ch <- val
			}
		}()
	}
	return dests
}

func main() {
	source := make(chan int)  // The input channel
	dests := Split(source, 5) // Retrieve 5 output channels

	go func() {
		for i := 1; i <= 10; i++ { // Send the number 1..10 to source and close it when we're done
			source <- i
		}

		close(source)
	}()

	var wg sync.WaitGroup // Use WaitGroup to wait until the output channels all close
	wg.Add(len(dests))

	for i, ch := range dests {
		go func(i int, d <-chan int) {
			defer wg.Done()

			for val := range d {
				fmt.Printf("#%d got %d\n", i, val)
			}
		}(i, ch)
	}

	wg.Wait()
}
