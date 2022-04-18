package throttle

import (
	"context"
	"fmt"
	"sync"
	"time"

	pkg_stp_types "github.com/4nt0nn/cn_go_patterns/pkg/stability_patterns"
)

func Throttle(effector pkg_stp_types.Effector, max uint, refill uint, d time.Duration) pkg_stp_types.Effector {
	var tokens = max
	var once sync.Once

	return func(ctx context.Context) (string, error) { // closure
		if ctx.Err() != nil {
			return "", ctx.Err()
		}

		once.Do(func() { // handles refill of tokens
			ticker := time.NewTicker(d)

			go func() {
				defer ticker.Stop()

				for {
					select {
					case <-ctx.Done():
						return
					case <-ticker.C:
						t := tokens + refill
						if t > max {
							t = max
						}
						tokens = t
					}
				}
			}()
		})

		if tokens <= 0 {
			return "", fmt.Errorf("to many calls")
		}

		tokens--

		return effector(ctx)
	}
}
