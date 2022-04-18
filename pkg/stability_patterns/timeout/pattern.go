package timeout

import (
	"context"

	pkg_stp_types "github.com/4nt0nn/cn_go_patterns/pkg/stability_patterns"
)

func Timeout(f pkg_stp_types.SlowFunction) pkg_stp_types.WithContext {
	return func(ctx context.Context, arg string) (string, error) {
		chres := make(chan string)
		cherr := make(chan error)

		go func() {
			res, err := f(arg)
			chres <- res
			cherr <- err
		}()

		select {
		case res := <-chres:
			return res, <-cherr
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
}

////// example use ///////

// func Slow(string) (string, error) {
// 	time.Sleep(2 * time.Second)
// 	return "HEY!", nil
// }

// func main() {
// 	ctx := context.Background()
// 	ctxt, cancel := context.WithTimeout(ctx, 1*time.Second)
// 	defer cancel()

// 	timeout := Timeout(Slow)
// 	res, err := timeout(ctxt, "some input")

// 	fmt.Println(res, err)
// }
