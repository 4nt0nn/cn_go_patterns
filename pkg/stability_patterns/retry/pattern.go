package retry

import (
	"context"
	"log"
	"time"

	pkg_stp_types "github.com/4nt0nn/cn_go_patterns/pkg/stability_patterns"
)

func Retry(effector pkg_stp_types.Effector, retires int, delay time.Duration) pkg_stp_types.Effector {
	return func(ctx context.Context) (string, error) {
		for r := 0; ; r++ { // condition omitted
			response, err := effector(ctx)
			if err == nil || r >= retires {
				return response, err
			}

			log.Printf("Attempt %d failed; retrying in %v", r+1, delay)

			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return "", ctx.Err()
			}
		}
	}
}

////// Example impl /////

// var count int

// func EmulateTransientError(ctx context.Context) (string, error) {
// 	count++

// 	if count <= 3 {
// 		return "intentional fail", errors.New("error")
// 	} else {
// 		return "success", nil
// 	}
// }

// func main() {
// 	r := Retry(EmulateTransientError, 5, 2*time.Second)

// 	res, err := r(context.Background())

// 	fmt.Println(res, err)
// }
