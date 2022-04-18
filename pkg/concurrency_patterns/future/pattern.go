package future

import (
	"context"
	"sync"
	"time"

	pkg_cnp_types "github.com/4nt0nn/cn_go_patterns/pkg/concurrency_patterns"
)

// InnerFuture is used internally to provide the concurrent functionality.
// In this example, it satisfies the Future interface, but could satisfy io.Reader by attaching a Read method.
type InnerFuture struct {
	once sync.Once
	wg sync.WaitGroup

	res string
	err error
	resCh <-chan string
	errCh <-chan error
}

// Result uses the attached structs once to call the passed function exactly once. Adds delta to the structs waitgroup which makes Future thread safe
// by blocking at f.wg.Wait() if any calls after the first is made until the channel reads are complete.
// Any values recieved by resCh and errCh will be assigned to the structs res and err field which is cached and returned.
func (f *InnerFuture) Result() (string, error) {
	f.once.Do(func() {
		f.wg.Add(1)
		defer f.wg.Done()
		f.res = <-f.resCh
		f.err = <-f.errCh
	})

	f.wg.Wait()

	return f.res, f.err
}

// SlowFunction is a wrapper around the core functionality that you want to run concurrently.
// It has the job of creating the result channels, running the core function in a goroutine, and 
// creating and returning the Future implementation (InnerFuture, in this example)
func SlowFunction(ctx context.Context) pkg_cnp_types.Future {
	resCh := make(chan string)
	errCh := make(chan error)

	go func() {
		select {
		case <- time.After(time.Second * 2):
			resCh <- "I slept for 2 seconds"
			errCh <- nil
		case <- ctx.Done():
			resCh <- ""
			errCh <- ctx.Err()
		}
	}()

	return &InnerFuture{resCh: resCh, errCh: errCh}
}

////// example use ///////

// This approach provides a reasonably good user experience.
// The programmer can create a Future and access it as they wish, and can even apply
// timeouts or deadlines with a context.

// func main() {
// 	ctx := context.Background()
// 	future := SlowFunction(ctx)

// 	res, err := future.Result()
// 	if err != nil {
// 		fmt.Println("error:", err)
// 		return
// 	}

// 	fmt.Println(res)
// }